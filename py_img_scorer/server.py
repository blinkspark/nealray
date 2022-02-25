from http.server import BaseHTTPRequestHandler, HTTPServer
from keras.models import load_model
import numpy as np
from PIL import Image

class Handler(BaseHTTPRequestHandler):
    model = load_model('model-dp2-06280.h5')

    def do_POST(self):
        data = self.rfile.read(int(self.headers['Content-Length']))
        img = Image.frombytes('RGB', (224, 224), data)
        img = np.array(img)
        img = img.reshape(1, 224, 224, 3)
        res = self.model.predict(img)
        res = np.argmax(res[0]) + 1
        print(res)
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(bytes(str(res), "utf-8"))


server = HTTPServer(('', 10088), Handler)
server.serve_forever()
