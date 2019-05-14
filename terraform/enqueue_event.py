import os
import json
import boto3
import urllib.parse
from typing import Dict


FORMENCODE_PREFIX = 'payload='


def challenge(payload: Dict) -> Dict:
    """Response to a Slack Event API challenge request"""
    return {
        "isBase64Encoded": False,
        "statusCode": 200,
        "headers": {'Content-Type': 'text/plain'},
        "body": payload['challenge']
    }


def enqueue_event(payload: Dict, app_key: str) -> Dict:
    """Enqueue an item in the queue"""
    queue_name = os.environ['queue_name']
    sqs = boto3.resource('sqs')
    queue = sqs.get_queue_by_name(QueueName=queue_name)
    msg = {"event": payload, "app_key": app_key}
    result = queue.send_message(MessageBody=json.dumps(msg), MessageGroupId='1666')
    return {
        "isBase64Encoded": False,
        "statusCode": 200,
        "headers": {'Content-Type': 'application/json'},
        "body": json.dumps({'queue': queue_name, 'result': result})
    }


def form_decode(payload: str) -> Dict:
    data = urllib.parse.parse_qs(payload)
    return json.loads(data['payload'][0])


def lambda_handler(event, context) -> Dict:
    if event['body'][:len(FORMENCODE_PREFIX)] == FORMENCODE_PREFIX:
        payload = form_decode(event['body'])
    else:
        payload = json.loads(event['body'])

    if 'challenge' in payload:
        return challenge(payload)

    app_key = event.get("pathParameters", {}).get("app_key", "unknown")
    return enqueue_event(payload, app_key)
