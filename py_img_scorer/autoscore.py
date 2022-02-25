#%%
import os
from keras.models import load_model
from PIL import Image
import numpy as np
import pandas as pd

workDir = 'o:/temp/pixiv/down'
model = load_model('./model-dp2-06280.h5')
data = pd.DataFrame()
for imgName in os.listdir(workDir):
    if not (imgName.endswith('.jpg') or imgName.endswith('.png'),
            imgName.endswith('.jpeg'), imgName.endswith('.JPG')
            or imgName.endswith('.PNG') or imgName.endswith('.JPEG')):
        continue

    imgPath = os.path.join(workDir, imgName)
    img = Image.open(imgPath)
    padding = max(img.size)
    paddingImg = Image.new('RGB', (padding, padding))
    paddingImg.paste(img, (int(
        (padding - img.size[0]) / 2), int((padding - img.size[1]) / 2)))
    img = paddingImg.resize((224, 224))
    img = np.array(img).reshape(1, 224, 224, 3)
    res = model.predict(img)[0]
    print(imgName, res)
    data.loc[imgName,'score'] = np.argmax(res) + 1

#%%
data.head()
#%%
data.to_csv('./result.csv', index_label='img')
