import argparse
import cv2 as cv
import numpy as np
import matplotlib.pyplot as plt
from datetime import datetime

face_jpg_file = r'face.jpg'
face_raw_file = r'face.raw'
# face_shape = (7680, 10240, 3)
face_shape = (15360, 20480, 3)
t = t_start = datetime.utcnow()

print('haveOpenCL', cv.ocl.haveOpenCL())
cv.ocl.setUseOpenCL(True)
print('useOpenCL', cv.ocl.useOpenCL())

def bench(msg):
    global t
    t0 = datetime.utcnow()
    dt = t0 - t
    t = t0
    print(t, msg, dt.total_seconds())

def load_image(filename, shape):
    """load_image load a given image. If the filename is empty, it loads a dummy image."""
    assert shape is not None, "an image shape must be given when loading raw files"
    face_memmap = np.memmap(filename, dtype=np.uint8, shape=shape)
    return face_memmap

def cvshow(img, **kwargs):
    img = img.clip(0, 255).astype('uint8')
    # CV stores colors as BGR; convert to RGB
    if img.ndim == 3:
        if img.shape[2] == 4:
            img = cv.cvtColor(img, cv.COLOR_BGRA2RGBA)
        else:
            img = cv.cvtColor(img, cv.COLOR_BGR2RGB)

    return plt.imshow(img, **kwargs)

p = argparse.ArgumentParser()
p.add_argument('-plot', '-p', default=False, action='store_true')
args = p.parse_args()

bench('start')

src = load_image(face_raw_file, face_shape)
src = cv.cvtColor(src, cv.COLOR_RGB2BGRA)
bench('rgb2bgra')

src_umat = cv.UMat(src)
bench('umat:mat')

srcTri = np.array([
    [0, 0],
    [src.shape[1] - 1, 0],
    [0, src.shape[0] - 1]]
).astype(np.float32)

dstTri = np.array([
    [0, src.shape[1]*0.33],
    [src.shape[1]*0.85, src.shape[0]*0.25],
    [src.shape[1]*0.15, src.shape[0]*0.7]]
).astype(np.float32)

warp_mat = cv.getAffineTransform(srcTri, dstTri)
warp_umat = cv.warpAffine(src_umat, warp_mat, (src.shape[1], src.shape[0]))
bench('warp_umat')
warp_dst = cv.UMat.get(warp_umat)
bench('warp_dst')

# Rotating the image after Warp
center = (warp_dst.shape[1]//2, warp_dst.shape[0]//2)
angle = -50
scale = 0.6

rot_mat = cv.getRotationMatrix2D(center, angle, scale)
warp_rotate_umat = cv.warpAffine(warp_umat, rot_mat, (warp_dst.shape[1], warp_dst.shape[0]))
bench('warp_rotate_umat')
# TODO: skip this step and use the raw data in memory
warp_rotate_dst = cv.UMat.get(warp_rotate_umat)
bench('warp_rotate_dst')

t = t_start
bench('end:transformation')

print(src.shape)

if args.plot:
    cvshow(warp_rotate_dst)
    plt.show()
    bench('plot')
