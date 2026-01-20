import os
from PIL import Image

INPUT_DIR = "actual_images"
OUTPUT_DIR = "resized_images"
TARGET_WIDTH = 600
THUMB_WIDTH = 200

def resize_images():
    if not os.path.exists(OUTPUT_DIR):
        os.makedirs(OUTPUT_DIR)

    files = [f for f in os.listdir(INPUT_DIR) if f.lower().endswith(('.jpg', '.jpeg', '.png'))]
    print(f"Found {len(files)} images.")

    for filename in files:
        input_path = os.path.join(INPUT_DIR, filename)
        
        # Open image
        try:
            with Image.open(input_path) as img:
                # 1. Main Image (600px width)
                w_percent = (TARGET_WIDTH / float(img.size[0]))
                h_size = int((float(img.size[1]) * float(w_percent)))
                img_resized = img.resize((TARGET_WIDTH, h_size), Image.Resampling.LANCZOS)
                
                output_path = os.path.join(OUTPUT_DIR, filename)
                img_resized.save(output_path, quality=85)
                print(f"Processed: {filename}")

                # 2. Thumbnail (200px width)
                # Naming convention: [filename]-sm.[ext] or just same name? 
                # Previous logic was likely [name]-sm.jpg. The frontend likely expects this or the DB does.
                # Let's check DB seed script naming convention to be sure. 
                # In seed_real_products.sql: '/api/media/products/v-rect-環遊世界-sm.jpg'
                
                name, ext = os.path.splitext(filename)
                thumb_filename = f"{name}-sm{ext}"
                
                t_percent = (THUMB_WIDTH / float(img.size[0]))
                t_h_size = int((float(img.size[1]) * float(t_percent)))
                thumb_resized = img.resize((THUMB_WIDTH, t_h_size), Image.Resampling.LANCZOS)
                
                thumb_path = os.path.join(OUTPUT_DIR, thumb_filename)
                thumb_resized.save(thumb_path, quality=85)
                print(f"Created thumbnail: {thumb_filename}")

        except Exception as e:
            print(f"Error processing {filename}: {e}")

if __name__ == "__main__":
    resize_images()
