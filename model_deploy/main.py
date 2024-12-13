from flask import Flask, request, jsonify, send_file
import numpy as np
import tensorflow as tf
from PIL import Image
import io

app = Flask(__name__)

# Load your TensorFlow model
model = tf.keras.models.load_model('nsfw_mobilenet.h5')

# Define the classes
classes = ['drawings', 'hentai', 'neutral', 'porn', 'sexy']

def prepare_image(image: Image.Image):
    """Preprocess the image to fit the model's input requirements."""
    image = image.resize((224, 224))
    image = image.convert('RGB')  # Ensure the image has 3 channels
    image_array = np.array(image)
    image_array = image_array / 255.0  # Normalize to 0-1
    image_array = np.expand_dims(image_array, axis=0)  # Add batch dimension
    return image_array

@app.route('/predict', methods=['POST'])
def predict():
    if 'file' not in request.files:
        return jsonify({'error': 'No file provided'}), 400
    file = request.files['file']
    if file:
        image = Image.open(io.BytesIO(file.read()))
        prepared_image = prepare_image(image)
        predictions = model.predict(prepared_image)
        predicted_class = classes[np.argmax(predictions)]
        confidence = np.max(predictions)

        # Determine label category
        label_category = 'Safe' if predicted_class in ['drawings', 'neutral'] else 'NSFW'

        # Logic to handle NSFW content
        if label_category == 'NSFW':
            return jsonify({'error': 'Content blocked due to NSFW image'}), 403
        
        # If safe, return the original image
        file.seek(0)  # Reset file pointer to the beginning
        return send_file(
            io.BytesIO(file.read()),
            mimetype='image/jpeg',
            as_attachment=True,
            attachment_filename=file.filename
        )
    else:
        return jsonify({'error': 'Invalid file'}), 400

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=8080)
