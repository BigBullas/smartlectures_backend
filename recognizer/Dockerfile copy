FROM python:3.11

WORKDIR /app

COPY ./requirements.txt ./


RUN apt update && apt install -y \
  tesseract-ocr-rus \
  ffmpeg \
  libsm6 \
  libxext6


RUN --mount=type=cache,target=/root/.cache/pip pip install -r requirements.txt

RUN --mount=type=cache,target=/root/.cache/pip pip install pix2text[multilingual]

COPY . .


EXPOSE 50051

CMD [ "python", "server.py" ]

