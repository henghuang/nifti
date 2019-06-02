import nibabel as nib
img = nib.load('MNI152.nii.gz')
print(img.shape)
data = img.get_fdata()
print(data[:,30:40,52])