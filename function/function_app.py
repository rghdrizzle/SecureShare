import logging
import azure.functions as func
from azure.identity import DefaultAzureCredential
from azure.storage.blob import BlobServiceClient, BlobClient, ContainerClient
import os

app = func.FunctionApp()

@app.timer_trigger(schedule="0 * * * *", arg_name="myTimer", run_on_startup=False,
              use_monitor=False) 
def delete_expired_files(myTimer: func.TimerRequest) -> None:
    if myTimer.past_due:
        logging.info('The timer is past due!')

    connect_str = os.getenv('AZURE_STORAGE_CONNECTION_STRING')
    blob_service_client = BlobServiceClient.from_connection_string(connect_str)
    blob_service_client.delete_blob(blob_service_client,"files",filename)
    logging.info('Python timer trigger function executed.')