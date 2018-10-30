import os
import json
import boto3


def challenge(payload):
    """Response to a Slack Event API challenge request"""
    return {
        "isBase64Encoded": False,
        "statusCode": 200,
        "headers": {'Content-Type': 'text/plain'},
        "body": payload['challenge']
    }


def enqueue(payload):
    """Enqueue an item in the queue"""
    queue_name = os.environ['queue_name']
    sqs = boto3.resource('sqs')
    queue = sqs.get_queue_by_name(QueueName=queue_name)
    result = queue.send_message(MessageBody=json.dumps(payload), MessageGroupId='1666')
    return {
        "isBase64Encoded": False,
        "statusCode": 200,
        "headers": {'Content-Type': 'application/json'},
        "body": json.dumps({'queue': queue_name, 'result': result})
    }


def lambda_handler(event, context):
    payload = json.loads(event['body'])
    if 'challenge' in payload:
        return challenge(payload)
    else:
        return enqueue(payload)
