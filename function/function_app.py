import logging
import azure.functions as func
from azure.identity import DefaultAzureCredential
from azure.storage.blob import BlobServiceClient, BlobClient, ContainerClient
from datetime import datetime, timezone
import re
import os

app = func.FunctionApp()

@app.timer_trigger(schedule="0 * * * *", arg_name="myTimer", run_on_startup=False,
              use_monitor=False) 
def delete_expired_files(myTimer: func.TimerRequest) -> None:
    if myTimer.past_due:
        logging.info('The timer is past due!')

    connect_str = os.getenv('AZURE_STORAGE_CONNECTION_STRING')
    blob_service_client = BlobServiceClient.from_connection_string(connect_str)
    container_client = blob_service_client.get_container_client(container="files")
    current_time = datetime.now(timezone.utc)
    blob_list = container_client.list_blobs()
    for blob in blob_list:
        blob_client = container_client.get_blob_client(blob.name)
        tags = blob_client.get_blob_tags()
        if "ExpiryTime" in tags:
            expiry_time_str = tags["ExpiryTime"]
            match = re.match(r"(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})", expiry_time_str)
            if match:
                expiry_time = datetime.strptime(match.group(1), "%Y-%m-%d %H:%M:%S").replace(tzinfo=timezone.utc)
                if current_time > expiry_time:
                    blob_client.delete_blob()
                    logging.info(f"{blob.name} - blob deleted successfully (expired at {expiry_time})")
                else:
                    logging.info(f"{blob.name} - not expired yet (expires at {expiry_time})")