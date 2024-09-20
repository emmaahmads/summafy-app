import PyPDF2
import boto3
import io
import json
import urllib3
import openai

print('Loading function')

# Initialize AWS clients
s3 = boto3.client('s3')
sns = boto3.client('sns')

# OpenAI API settings
openai.api_key = "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"

def generate_summary(book_name, book_content):
    prompt = f"Summarize the book '{book_name}' in brief. The book is about {book_content}."
    response = openai.Completion.create(
        engine="text-davinci-002",  # You can use a different engine if you prefer
        prompt=prompt,
        temperature=0.7,
        max_tokens=200,
        stop=["\n"],
    )
    return response.choices[0].text.strip()
    
def pdf_to_text(pdf_file):
    """Extract text from a PDF file."""
    pdf_reader = PyPDF2.PdfReader(pdf_file)
    text = ""

    for page in pdf_reader.pages:
        text += page.extract_text()

    return text

def lambda_handler(event, context):
      # Extract event data
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = event['Records'][0]['s3']['object']['key']
    event_name = event['Records'][0]['eventName']

    # Get object content
    response = s3.get_object(Bucket=bucket, Key=key)
    object_content = response['Body'].read()
    pdf_file = io.BytesIO(object_content)
    text = pdf_to_text(pdf_file)
   # text = text[:15000]

    # Prepare OpenAI API request
    headers = {
        "Authorization": f"Bearer {OPENAI_API_KEY}",
        "Content-Type": "application/json"
    }
    data = {
        "model": "text-davinci-002",
        "prompt": f"Write a short summary of the following text: {text}",
        temperature=0.7,
        max_tokens=200,
        stop=["\n"],
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