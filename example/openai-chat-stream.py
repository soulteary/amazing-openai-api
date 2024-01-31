from openai import OpenAI

client = OpenAI(
    api_key="your-key-or-input-something-as-you-like",
    base_url="http://127.0.0.1:8080/v1"
)

stream = client.chat.completions.create(
    model="gpt-4",
    messages=[{"role": "user", "content": "Write a romantic poem and talk about League of Legends"}],
    stream=True,
)
for chunk in stream:
    print(chunk.choices[0].delta.content or "", end="")