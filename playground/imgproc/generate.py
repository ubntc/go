import cv2 as cv
from scipy import misc

def generate_image() -> bool:
    face = misc.face()
    face = cv.vconcat([face] * 20)
    face = cv.hconcat([face] * 20)
    cv.imwrite('face.jpg', face)
    face.tofile('face.raw')
    print(f"""
    Saved raw image with Image Shape: {face.shape}
    ATTENTION: You must use this shape as `face_shape` in the script!
    """)

if __name__ == '__main__':
    generate_image()
