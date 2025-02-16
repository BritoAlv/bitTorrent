from flask import Flask, send_file

app = Flask(__name__)

@app.route('/', methods = ['GET'])
def home():
    return send_file("./static/index.html")

@app.route('/index', methods = ['GET'])
def index():
    return home()

if __name__ == '__main__':
    app.run(debug=True, port=9090)