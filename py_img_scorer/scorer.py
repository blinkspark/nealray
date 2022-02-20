import shutil
import tkinter as tk
from tkinter import ttk
from tkinter import filedialog
from PIL import Image, ImageTk
import os, sys, re
from os import path
import pandas as pd


####################### ImageGenerator #############################
class ImageGenerator:

    def __init__(self, workingDir, targetImgHeight=1000):
        self.workingDir = workingDir
        self.index = 0
        self.imgNameList = list(self.img_gen(workingDir))
        self.currentImg = None
        self.nextImg = None
        self.prevImg = None
        self.targetImgHeight = targetImgHeight

        if len(self.imgNameList) > 0:
            self.init_img()

    def init_img(self):
        self.currentImg = self.img_tk(self.index)
        self.nextImg = self.img_tk(self.next_index(self.index))
        self.prevImg = self.img_tk(self.prev_index(self.index))

    @staticmethod
    def img_gen(wrorkingDir):
        for f in os.listdir(wrorkingDir):
            if f.endswith('.jpg') or f.endswith('.png') or f.endswith(
                    '.jpeg') or f.endswith('.JPG') or f.endswith(
                        '.PNG') or f.endswith('.JPEG'):
                yield f

    @staticmethod
    def resize_img(img, width=-1, height=-1, inter=Image.BICUBIC):
        if width < 0 and height < 0:
            return img
        elif width < 0:
            img_size = img.size
            resize_ratio = height / img_size[1]
            new_size = (round(img_size[0] * resize_ratio),
                        round(img_size[1] * resize_ratio))
            return img.resize(new_size, inter)
        elif height < 0:
            img_size = img.size
            resize_ratio = width / img_size[0]
            new_size = (round(img_size[0] * resize_ratio),
                        round(img_size[1] * resize_ratio))
            return img.resize(new_size, inter)

    def img_tk(self, index):
        img = Image.open(path.join(self.workingDir, self.imgNameList[index]))
        img = self.resize_img(img, height=self.targetImgHeight)
        img = ImageTk.PhotoImage(img)
        return img

    def set_working_dir(self, workingDir):
        self.workingDir = workingDir
        self.index = 0
        self.imgNameList = list(self.img_gen(workingDir))
        self.init_img()

    def next_index(self, index):
        return (index + 1) % len(self.imgNameList)

    def prev_index(self, index):
        return (index - 1) % len(self.imgNameList)

    def next(self):
        self.index = self.next_index(self.index)
        self.prevImg = self.currentImg
        self.currentImg = self.nextImg
        self.nextImg = self.img_tk(self.next_index(self.index))
        return self.currentImg

    def prev(self):
        self.index = self.prev_index(self.index)
        self.nextImg = self.currentImg
        self.currentImg = self.prevImg
        self.prevImg = self.img_tk(self.prev_index(self.index))
        return self.currentImg

    def set_index(self, index):
        self.index = index
        self.init_img()

    def get_img(self):
        return self.currentImg

    def dedupe(self):
        for imgName in self.imgNameList:
            if re.search(r'\(\d+\)', imgName):
                os.remove(path.join(self.workingDir, imgName))
        self.imgNameList = list(self.img_gen(self.workingDir))
        self.index = 0
        self.init_img()


############################## APP ###################################
class App(tk.Tk):

    def __init__(self):
        super().__init__()
        self.title("Image Scorer")

        self.workingDir = sys.argv[1] if (
            len(sys.argv) > 1 and os.path.exists(sys.argv[1])) else '.'

        self.dataCsvPath = path.join(self.workingDir, 'data.csv')
        self.dataFrame = pd.read_csv(
            self.dataCsvPath, index_col='id') if path.exists(
                self.dataCsvPath) else pd.DataFrame(columns=['img', 'score'])

        self.imgGenerator = ImageGenerator(self.workingDir)
        self.currImg = self.imgGenerator.get_img()

        self.mainFrame = self.main_frame(self)
        self.mainMenu = self.main_menu(self)
        self.imageFrame = self.img_frame(self.mainFrame)
        self.imgLabel = self.img_label(self.imageFrame)
        self.nextBtn = self.next_btn(self.imageFrame, self.next_action)
        self.prevBtn = self.back_btn(self.imageFrame, self.prev_action)
        self.prevBtn.configure(command=self.prev_action)
        self.infoFrame = self.info_frame(self.mainFrame)
        self.indexLabel = self.index_label(self.infoFrame)
        self.pathLabel = self.path_label(self.infoFrame)
        self.scoreLabel = self.score_label(self.infoFrame)

        self.next_unscored()
        # self.update_all()

        self.bind('<Right>', lambda e: self.next_action())
        self.bind('<Left>', lambda e: self.prev_action())
        self.bind(
            '<Home>', lambda e:
            (self.imgGenerator.set_index(0), self.update_all()))
        self.bind('<Escape>', lambda e: self.quit())

        self.bind('1', lambda e: self.set_score(1))
        self.bind('2', lambda e: self.set_score(2))
        self.bind('3', lambda e: self.set_score(3))
        self.bind('4', lambda e: self.set_score(4))
        self.bind('5', lambda e: self.set_score(5))
        self.bind('s', lambda e: self.save_csv_action())

    @staticmethod
    def main_frame(parent):
        mainFrame = ttk.Frame(parent)
        mainFrame.pack(fill=tk.BOTH, expand=True)
        return mainFrame

    @staticmethod
    def img_frame(parent):
        imgFrame = ttk.Frame(parent)
        imgFrame.grid(row=0, column=0, columnspan=3, sticky=tk.NSEW)
        return imgFrame

    @staticmethod
    def img_label(parent):
        imgLabel = ttk.Label(parent)
        imgLabel.grid(row=1, column=1, sticky=tk.NSEW)
        return imgLabel

    @staticmethod
    def next_btn(parent, command):
        nextBtn = ttk.Button(parent, text='>>', command=command)
        nextBtn.grid(row=1, column=2, sticky=tk.NSEW)
        return nextBtn

    @staticmethod
    def back_btn(parent, command):
        backBtn = ttk.Button(parent, text='<<', command=command)
        backBtn.grid(row=1, column=0, sticky=tk.NSEW)
        return backBtn

    @staticmethod
    def main_menu(root):
        mainMenu = tk.Menu(root)

        fileMenu = tk.Menu(mainMenu, tearoff=0)
        fileMenu.add_command(label='Open', command=root.file_open_action)
        fileMenu.add_command(label='Save', command=root.save_csv_action)
        fileMenu.add_separator()
        fileMenu.add_command(label='Dedupe', command=root.dedupe_action)
        fileMenu.add_command(label='Merge...', command=root.merge_action)
        fileMenu.add_separator()
        fileMenu.add_command(label='Exit', command=lambda: root.quit())
        mainMenu.add_cascade(label='File', menu=fileMenu)

        root.config(menu=mainMenu)
        return mainMenu

    @staticmethod
    def info_frame(parent):
        infoFrame = ttk.Frame(parent)
        # set pack direction
        infoFrame.grid(row=2, column=0, columnspan=3)
        return infoFrame

    @staticmethod
    def index_label(parent):
        indexLabel = ttk.Label(parent, text='0/0', padding=(5, 0))
        indexLabel.pack(side=tk.LEFT)
        return indexLabel

    @staticmethod
    def path_label(parent):
        pathLabel = ttk.Label(parent, text='', padding=(5, 0))
        pathLabel.pack(side=tk.LEFT)
        return pathLabel

    @staticmethod
    def score_label(parent):
        scoreLabel = ttk.Label(parent, text='0', padding=(5, 0))
        scoreLabel.pack(side=tk.LEFT)
        return scoreLabel

    def update_img(self):
        self.currImg = self.imgGenerator.get_img()
        if self.currImg:
            self.imgLabel.config(image=self.currImg)

    def update_index(self):
        self.indexLabel.configure(
            text=
            f'{self.imgGenerator.index + 1}/{len(self.imgGenerator.imgNameList)}'
        )

    def update_path(self):
        self.pathLabel.configure(text=path.join(
            self.workingDir, self.imgGenerator.imgNameList[
                self.imgGenerator.index]))

    def update_score(self):
        score = self.dataFrame.loc[self.dataFrame['img'] ==
                                   self.imgGenerator.imgNameList[
                                       self.imgGenerator.index]]['score']
        self.scoreLabel.configure(
            text=score.values[0] if score.size > 0 else 0)

    def update_all(self):
        self.update_img()
        self.update_index()
        self.update_path()
        self.update_score()

    def file_open_action(self):
        newDir = filedialog.askdirectory(initialdir=self.workingDir,
                                         mustexist=True)
        if newDir == '':
            return
        self.workingDir = newDir
        self.imgGenerator.set_working_dir(self.workingDir)
        self.currImg = self.imgGenerator.get_img()
        self.next_unscored()
        # self.update_all()

    def next_action(self):
        self.currImg = self.imgGenerator.next()
        self.update_all()

    def prev_action(self):
        self.currImg = self.imgGenerator.prev()
        self.update_all()

    def dedupe_action(self):
        self.imgGenerator.dedupe()
        self.update_all()

    def set_score(self, score):
        scoreData = pd.DataFrame(
            {
                'img': self.imgGenerator.imgNameList[self.imgGenerator.index],
                'score': score
            },
            index=[self.imgGenerator.index])
        self.dataFrame = pd.concat([self.dataFrame, scoreData],
                                   ignore_index=True)
        self.dedupe_data()
        self.imgGenerator.next()
        self.update_all()

    def next_unscored(self):
        for (i, img) in enumerate(self.imgGenerator.imgNameList):
            self.dataFrame['img'].isin([img]).any()
            if img not in self.dataFrame['img'].values:
                self.imgGenerator.set_index(i)
                self.currImg = self.imgGenerator.get_img()
                self.update_all()
                break
        self.update_all()

    def dedupe_data(self):
        self.dataFrame.drop_duplicates(subset='img', inplace=True, keep='last')
        self.dataFrame.reset_index(drop=True, inplace=True)

    def save_csv_action(self):
        self.dataFrame.to_csv(self.dataCsvPath, index_label='id')

    def merge_action(self):
        newDir = filedialog.askdirectory(initialdir=self.workingDir, mustexist=True)
        if newDir == '':
            return
        newDataFrame = pd.read_csv(path.join(newDir, 'data.csv'), index_col='id')
        self.dataFrame = pd.concat([self.dataFrame, newDataFrame], ignore_index=True)
        self.dedupe_data()
        self.save_csv_action()
        
        imgGen = ImageGenerator(newDir)
        for imgName in imgGen.imgNameList:
            imgPath = path.join(newDir, imgName)
            imgTargetPath = path.join(self.workingDir, imgName)
            if not path.exists(imgTargetPath):
                shutil.copy(imgPath, imgTargetPath)

        self.imgGenerator.set_working_dir(self.workingDir)
        self.imgGenerator.set_index(0)
        self.update_all()


app = App()
app.mainloop()
