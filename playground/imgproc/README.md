# Image Processing Scripts

This project contains scripts and research results for some image processing needs I encountered here and there.

All code is provided as is and may not represent any state of the art in computer vision. Feel free to open an issue here to teach me how I can improve my approaches.

# Ressources

The following articles and docs were used to develop the solutions:

* https://learnopencv.com/opencv-transparent-api/
* https://github.com/jupyter/notebook/issues/3935
* https://www.google.com/search?q=opencv+affine+transformation
* https://docs.opencv.org/3.4/d4/d61/tutorial_warp_affine.html
* https://note.nkmk.me/en/python-opencv-hconcat-vconcat-np-tile/

# Script 1: Affine Transformation

## Requirements

* fast processing in GPU
* support for super-large image sizes (>700 MB raw)
* use common technologies that are widely available

## Solution

I chose to use OpenCV, which provides GPU support via its Transparent API. Also it is widely available and implementations can be adopted across languages.

This solution uses Python, which was the easiest for rapid prototyping for me.

## Dependencies
You mainly need OpenCV. For this script I also use SciPy for getting a sample image and some common iPython Notebook tooling for visualizing the results.
```
pip install opencv-contrib-python
pip install numpy scipy matplotlib ipython jupyter
```

## Usage

Generate a ⚠️**1GB**⚠️ raw image for testing.
```
python3 generate.py
```

Run a transformation on the raw image.
```
python3 transform.py     # transform and exit
python3 transform.py -p  # transform and plot
```

Remove generated files.
```
make clean
````

Also see [Makefile](Makefile)

## Learnings

Use `cv.UMat` to make use of OpenCV TAPI, which is supposed to run on the GPU.

```python
src_umat = cv.UMat(src)
warp_umat = cv.warpAffine(src_umat, ...)
```

The `UMat` does not have a `shape`, so you need to compute all shape logic separately.
The script still uses

```python
warp_dst = cv.UMat.get(warp_umat)
```

to get the "image" back from the "matrix" in the GPU.

This is expensive! You need to do the math to compute the target shape `(w, h, colors)` after each step,
since the shape is usually needed to set up the next transformation.

But theoretically it should be possible to prepare all steps and shapes before running the actual transformations, so that only the computation-heavy transformations are left to be run on the GPU.