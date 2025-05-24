from pix2text import Pix2Text
from text_recognizer import TextRecognizerRus
from PIL import Image
from io import BytesIO
import argparse
import re


class LatexOCR:
    def __init__(self):
        self.p2t = Pix2Text(analyzer_config=dict(model_name='mfd'))
        self.text_rec = TextRecognizerRus()

    def convert(self, data: bytes) -> str:
        bytes_stream = BytesIO(data)
        image = Image.open(bytes_stream)
        outs = self.p2t.recognize(image, resized_shape=608, return_text=False)

        line_number = 0
        markdown = ''
        for element in outs:
            text = element['text']
            if element['type'] == 'text':
                top_left = element['position'][0]
                bottom_right = element['position'][2]
                coords = tuple(list(top_left) + list(bottom_right))
                cropped = image.crop(coords)
                text = self.text_rec.recognize(cropped).strip()
            
            elif element['type'] == 'embedding':
                text = ' `$$' + text + '$$`'
            elif element['type'] == 'isolated':
                text = '\n' + '```KaTeX\n' + text + '\n```' + '\n'

            if element['line_number'] == line_number:
                markdown += text
            else:
                markdown += '\n' + text
                line_number = element['line_number']

        # markdown = re.sub(r'\\tag{(?P<number>\d+)}', r'(\g<number>)', markdown)
        markdown = re.sub(r'\\tag{\d+.?\d.}', '', markdown)

        return markdown
    
    def convert_text(self, data: bytes) -> str:
        # bytes_stream = BytesIO(data)
        # image = Image.open(bytes_stream)
        # text = self.text_rec.recognize(image).strip()
        # return text
        return self.convert(data)
    
    def convert_formula(self, data: bytes) -> str:
        # bytes_stream = BytesIO(data)
        # image = Image.open(bytes_stream)
        # formula = self.p2t.recognize_formula(image)
        # text = '`$$' + formula + '$$`'
        # return text
        return self.convert(data)
    
    def convert_mixed(self, data: bytes) -> str:
        return self.convert(data)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('image_path', type=str, help='Path to the image file')
    args = parser.parse_args()

    with open(args.image_path, 'rb') as img:
        img_bytes = img.read()

    latex_ocr = LatexOCR()
    markdown = latex_ocr.convert(img_bytes)
    
    with open('out.md', 'w') as out:
        out.write(markdown)
