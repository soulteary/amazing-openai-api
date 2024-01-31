from openai import OpenAI

client = OpenAI(
    api_key="your-key-or-input-something-as-you-like",
    base_url="http://127.0.0.1:8080/v1"
)

chat_completion = client.chat.completions.create(
    messages=[
        {
            "role": "user",
            "content": "Say this is a test",
        }
    ],
    model="gpt-3.5-turbo",
)

print(chat_completion)