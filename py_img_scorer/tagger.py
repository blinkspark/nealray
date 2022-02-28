import sys
import tkinter as tk, tkinter.ttk as ttk, tkinter.filedialog as filedialog
import numpy as np, pandas as pd
from keras.models import load_model
import os, os.path as path
from PIL import Image, ImageTk


class APP(tk.Tk):

    def __init__(self, workingDir='.'):
        super().__init__()

        self.workingDir = workingDir
        self.tagNames = ['R18', 'Quality']

        self.title("Tagger")
        self.mode = tk.StringVar(value=self.tagNames[0])

        self.index = -1
        self.load_data()
        self.imgList = list(self.img_list_generator())

        ttk.Style().configure("TRadiobutton", font=(18))
        self.geometry("1200x900")

        self.main_menu()
        self.main_frame()
        self.mode_select_frame()
        self.img_label()
        self.prev_btn()
        self.next_btn()
        self.info_frame()
        self.goto_frame()

        self.set_index(0)
        # self.update_all()

        self.bind("<Configure>", self.on_resize)
        self.bind("<Escape>", lambda e: self.on_destroy_action())
        self.bind("<Left>", lambda e: self.prev_btn_action())
        self.bind("<Right>", lambda e: self.next_btn_action())
        self.bind("<Home>", lambda e: self.home_action())

        self.bind('1', lambda e: self.tag_action(1))
        self.bind('2', lambda e: self.tag_action(2))
        self.bind('`', lambda e: self.tag_action(np.nan))

        self.protocol("WM_DELETE_WINDOW", self.on_destroy_action)

    ##################### actions #####################
    def on_destroy_action(self):
        try:
            self.data.to_csv(path.join(self.workingDir, 'data.csv'),
                             index_label='img')
        except Exception as e:
            print(e)
        self.destroy()

    def on_resize(self, event):
        self.update_all()

    def prev_btn_action(self):
        self.set_index((self.index - 1) % len(self.imgList))

    def next_btn_action(self):
        self.set_index((self.index + 1) % len(self.imgList))

    def home_action(self):
        self.set_index(0)

    def goto_action(self):
        self.set_index(int(self.gotoEntry.get()) - 1)
        self.focus_set()

    def tag_action(self, score):
        if self.focus_get() == self.gotoEntry:
            return
        tagName = self.mode.get()
        self.data.loc[self.imgList[self.index], tagName] = float(score - 1)
        self.next_btn_action()

    def reset_action(self):
        self.data = pd.DataFrame()
        dataPath = path.join(self.workingDir, 'data.csv')
        if path.exists(dataPath):
            os.remove(dataPath)
        self.set_index(0)

    ##################### update #####################
    def update_all(self):
        self.update_img()
        self.update_index_label()
        self.update_name_label()
        self.update_tag_labels()

    def update_img(self):
        winSize = self.mainFrame.grid_bbox(row=1, column=1)
        winSize = (winSize[2] - winSize[0], winSize[3] - winSize[1])
        if winSize[0] <= 0 or winSize[1] <= 0:
            return
        img = self.currentImg
        newSize = self.get_new_img_size(winSize, img.size)
        if min(newSize) <= 0:
            return
        img = img.resize(newSize, Image.BICUBIC)
        img = ImageTk.PhotoImage(img)
        self.imgLabel.img = img
        self.imgLabel.configure(image=img)

    def get_new_img_size(self, winSize, imgSize):
        ratio = winSize[0] / imgSize[0] if imgSize[0] > imgSize[
            1] else winSize[1] / imgSize[1]
        newSize = (int(imgSize[0] * ratio), int(imgSize[1] * ratio))
        if newSize[0] > winSize[0]:
            ratio = winSize[0] / newSize[0]
            newSize = (int(newSize[0] * ratio), int(newSize[1] * ratio))
        if newSize[1] > winSize[1]:
            ratio = winSize[1] / newSize[1]
            newSize = (int(newSize[0] * ratio), int(newSize[1] * ratio))

        return newSize

    def load_img(self):
        img = Image.open(path.join(self.workingDir, self.imgList[self.index]))
        self.currentImg = img

    def set_index(self, i):
        if i == self.index:
            return
        self.index = i
        self.load_img()
        self.update_all()

    def update_index_label(self):
        self.indexLabel.configure(text=f'{self.index + 1}/{len(self.imgList)}')

    def update_name_label(self):
        self.nameLabel.configure(text=self.imgList[self.index])

    def update_tag_labels(self):
        for label in self.tagLabels:
            try:
                tag = label.tag
                label.configure(
                    text=
                    f'{tag}: {self.data.loc[self.imgList[self.index], tag]}')
            except Exception as e:
                label.configure(text=f'{tag}: {np.nan}')

    ##################### data #####################
    def load_data(self):
        dataPath = path.join(self.workingDir, 'data.csv')
        if path.exists(dataPath):
            self.data = pd.read_csv(dataPath, index_col='img')
        else:
            self.data = pd.DataFrame()

    def img_list_generator(self):
        for f in os.listdir(self.workingDir):
            if not (f.endswith('.jpg') or f.endswith('.png')
                    or f.endswith('.jpeg') or f.endswith('.JPG')
                    or f.endswith('.PNG') or f.endswith('.JPEG')):
                continue
            yield f

    #################### widgets ####################
    def main_menu(self):
        mainMenu = tk.Menu(self)

        fileMenu = tk.Menu(mainMenu, tearoff=0)
        fileMenu.add_command(label="Open")
        fileMenu.add_command(label="Save")
        fileMenu.add_command(label="Reset", command=self.reset_action)
        fileMenu.add_separator()
        fileMenu.add_command(label="Exit", command=self.on_destroy_action)
        mainMenu.add_cascade(label="File", menu=fileMenu)

        self.config(menu=mainMenu)
        self.mainMenu = mainMenu

    def main_frame(self):
        self.mainFrame = ttk.Frame(self)
        self.mainFrame.grid(row=0, column=0, sticky=tk.NSEW)
        self.rowconfigure(0, weight=1)
        self.columnconfigure(0, weight=1)

    def mode_select_frame(self):
        msf = ttk.Frame(self.mainFrame)
        for text in self.tagNames:
            rb = ttk.Radiobutton(msf,
                                 text=text,
                                 variable=self.mode,
                                 value=text)
            rb.pack(side=tk.LEFT, padx=5, pady=5)
        msf.grid(row=0, column=0, columnspan=3)
        self.modeSelectFrame = msf

    def img_label(self):
        imgLabel = ttk.Label(self.mainFrame, text="Image")
        imgLabel.grid(row=1, column=1)
        self.mainFrame.rowconfigure(1, weight=1)
        self.mainFrame.columnconfigure(1, weight=1)
        self.imgLabel = imgLabel

    def prev_btn(self):
        prevBtn = ttk.Button(self.mainFrame,
                             text="<<",
                             command=self.prev_btn_action)
        prevBtn.grid(row=1, column=0, sticky=tk.NSEW)
        self.prevBtn = prevBtn

    def next_btn(self):
        nextBtn = ttk.Button(self.mainFrame,
                             text=">>",
                             command=self.next_btn_action)
        nextBtn.grid(row=1, column=2, sticky=tk.NSEW)
        self.nextBtn = nextBtn

    def info_frame(self):
        infoFrame = ttk.Frame(self.mainFrame)

        indexLabel = ttk.Label(infoFrame, text="Index")
        indexLabel.pack(side=tk.LEFT, padx=5, pady=5)
        self.indexLabel = indexLabel

        nameLabel = ttk.Label(infoFrame, text="Path")
        nameLabel.pack(side=tk.LEFT, padx=5, pady=5)
        self.nameLabel = nameLabel

        self.tagLabels = []
        for labelName in self.tagNames:
            label = ttk.Label(infoFrame, text=labelName)
            label.tag = labelName
            label.pack(side=tk.LEFT, padx=5, pady=5)
            self.tagLabels.append(label)

        infoFrame.grid(row=2, column=0, columnspan=3)
        self.infoFrame = infoFrame

    def goto_frame(self):
        gotoFrame = ttk.Frame(self.mainFrame)

        gotoLabel = ttk.Label(gotoFrame, text="goto:")
        gotoLabel.pack(side=tk.LEFT, padx=5, pady=5)

        gotoEntry = ttk.Entry(gotoFrame)
        gotoEntry.pack(side=tk.LEFT, padx=5, pady=5)
        gotoEntry.bind("<Return>", lambda e: self.goto_action())
        self.gotoEntry = gotoEntry

        gotoFrame.grid(row=3, column=0, columnspan=3)
        self.gotoFrame = gotoFrame


workingDir = sys.argv[1] if len(sys.argv) > 1 else '.'
app = APP(workingDir=workingDir)
app.mainloop()
