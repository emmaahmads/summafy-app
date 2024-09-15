import json
import boto3
import urllib3

print('Loading function')

# Initialize AWS clients
s3 = boto3.client('s3')
sns = boto3.client('sns')

# OpenAI API settings
OPENAI_API_KEY = "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"
OPENAI_API_URL = "https://api.openai.com/v1/completions"

def lambda_handler(event, context):
    # Extract event data
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = event['Records'][0]['s3']['object']['key']
    event_name = event['Records'][0]['eventName']

    # Get object content
    response = s3.get_object(Bucket=bucket, Key=key)
    object_content = response['Body'].read().decode('utf-8')

    # Prepare OpenAI API request
    headers = {
        "Authorization": f"Bearer {OPENAI_API_KEY}",
        "Content-Type": "application/json"
    }
    data = {
        "model": "gpt-3.5-turbo-instruct",
        "prompt": f"Write a short summary of the following text: {object_content}",
        "max_tokens": 100,
        "temperature": 0.5
    }

    # Call OpenAI API
    http = urllib3.PoolManager()
    ai_response = http.request('POST', OPENAI_API_URL, headers=headers, body=json.dumps(data))

    # Parse JSON response
    response_data = json.loads(ai_response.data.decode('utf-8'))
    print(response_data)  # Print the response data for debugging

    # Check if the 'choices' key exists
    if 'choices' in response_data:
    # Check if the 'choices' list is not empty
        if len(response_data['choices']) > 0:
            # Check if the first choice has a 'text' key
            if 'text' in response_data['choices'][0]:
                summary = response_data["choices"][0]["text"]
            else:
                print("Error: The first choice does not have a 'text' key.")
        else:
            print("Error: The 'choices' list is empty.")
    else:
        print("Error: The 'choices' key does not exist.")

    # Prepare SNS message
    sns_message = f"This Email Represent a File Status has been Changed in One of Your Bucket\n\nBUCKET NAME: {bucket}\n\nFILE NAME: {key}\n\nSUMMARY: {summary}"

    # Send SNS message
    sns.publish(
        TopicArn="arn:aws:sns:us-east-1:975050042748:s3_newfile_send_email",
        Message=sns_message,
        Subject="File changes summary"
    )

    return {
        'statusCode': 200,
        'statusMessage': 'OK'
    }

""" example response:
{'id': 'cmpl-A5kGAzDy1thrG3kaRHNZ0HvZTTgZT', 'object': 'text_completion', 'created': 1725932826, 'model': 'gpt-3.5-turbo-instruct', 'choices': [{'text': '\r\n\r\nThe article discusses the trend of people delaying having children until their late 30s or 40s, but warns that as a woman ages, her fertility declines and the chances of getting pregnant naturally become unlikely after age 45. This is due to a decrease in the number of eggs and an increased likelihood of abnormal chromosomes. The article also mentions the importance of considering an infertility evaluation, especially for those over 35 or with existing fertility issues. Getting pregnant later in life can also pose risks such', 'index': 0, 'logprobs': None, 'finish_reason': 'length'}], 'usage': {'prompt_tokens': 524, 'completion_tokens': 100, 'total_tokens': 624}}
fertility.txt """