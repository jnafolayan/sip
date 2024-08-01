# Sip
Image compression using wavelet transforms. Play with it here: https://jnafolayan.github.io/sip/web/.

**Wavelet families implemented**
- Haar
- Cohen-Daubechies-Feauveau: this is particularly more CPU intensive, so expect longer compression times for large decomposition levels.

## Usage
If you'd like to run it locally, follow these steps:

```sh
git clone https://github.com/jnafolayan/sip
cd sip

# Compress an image
# wavelet types: haar, cdf97
# level (levels of decomposition): 1, 2, 3, ...
# threshold: 0, 1, 2, ...
go run main.go compress PATH_TO_FILE -wavelet WAVELET_TYPE -level LEVEL -output OUTPUT_FILE -threshold THRESHOLD
```