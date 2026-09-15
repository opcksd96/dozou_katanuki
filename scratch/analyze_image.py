from PIL import Image
import numpy as np

img = Image.open(r"C:\Users\Const\.gemini\antigravity-ide\brain\447b7f46-7ce4-4e71-bf7f-5b698f69daf1\.user_uploaded\media_1789169176046.png")
print("Image size:", img.size)

# Crop one of the boxes
# The boxes are roughly from y=150 to y=350, x=50 to x=400
arr = np.array(img)
print("Image shape:", arr.shape)

# Let's inspect colors along a vertical line or inside a box
# Box 1 roughly at y=160..180
# Print unique colors in a region
box = arr[140:190, 50:400]
print("Unique colors in box 1 (first 20):")
unique_colors = np.unique(box.reshape(-1, arr.shape[2]), axis=0)
print(unique_colors[:20])
