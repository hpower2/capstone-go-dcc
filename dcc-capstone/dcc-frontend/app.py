import json
from flask import Flask, jsonify, render_template, request, url_for, redirect, session
import requests
from flask_cors import CORS


app = Flask(__name__, static_url_path='/static')
CORS(app)

@app.route('/')
def index():
    trainName = requests.get("http://localhost:8000/train").json()
    return render_template('index.html', trainName = trainName)


@app.route('/api/test', methods=['POST'])
def test():
    data = request.json
    print(data)
    return {"data" : data}

@app.route('/command', methods=['POST'])
def command():
    data = request.json
    data = json.dumps(data)
    print(data)
    requests.post(f'http://localhost:8000/command', data=data)
    return {"None" : "None"}

if __name__ == "__main__":
    app.run(debug=True)


