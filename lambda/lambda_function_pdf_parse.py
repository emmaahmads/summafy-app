import PyPDF2

print('Loading function')

def pdf_to_text(pdf_file):
    """Extract text from a PDF file."""
    pdf_reader = PyPDF2.PdfReader(pdf_file)
    text = ""

    for page in pdf_reader.pages:
        text += page.extract_text()

    return text

def main():
    with open('test.pdf', 'rb') as file:
        print(pdf_to_text(file))


if __name__ == '__main__':
    main()
