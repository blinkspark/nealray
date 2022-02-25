import os, sys, re, shutil
import tkinter as tk
from tkinter import ttk, filedialog
from PIL import Image, ImageTk
from os import path
import pandas as pd
import configparser as cp
import numpy as np


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

    def get_img_name(self):
        return self.imgNameList[self.index]

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
        ttk.Style().configure('TCheckbutton', font=(18))
        self.lastTagRow = -1
        self.lastTagCol = -1
        self.cfg = cp.ConfigParser()
        self.cfg.read('config.ini')

        self.workingDir = sys.argv[1] if (
            len(sys.argv) > 1 and os.path.exists(sys.argv[1])) else '.'

        self.dataCsvPath = path.join(self.workingDir, 'data.csv')
        self.scoreCsvPath = path.join(self.workingDir, 'score.csv')
        self.dataFrame = self.load_data()
        self.scoreData = self.load_score_data()

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
        self.tagFrame = self.tag_frame(self.mainFrame)
        self.tagCheckBoxes = self.tag_checkboxes(self.tagFrame)
        self.tagBtn = self.tag_btn(self.tagFrame, self.tag_action)

        # self.next_untaged()
        # self.update_all()
        self.next_unscored()

        self.bind('<Right>', lambda e: self.next_action())
        self.bind('<Left>', lambda e: self.prev_action())
        self.bind(
            '<Home>', lambda e:
            (self.imgGenerator.set_index(0), self.update_all()))
        self.bind('<Escape>', lambda e: self.quit())

        #scores binding
        self.bind('1', lambda e: self.set_score(1))
        self.bind('2', lambda e: self.set_score(2))
        self.bind('3', lambda e: self.set_score(3))

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
        fileMenu.add_command(label='Reset', command=root.reset_action)
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
    def tag_frame(parent):
        tagFrame = ttk.Frame(parent)
        tagFrame.grid(row=3, column=0, columnspan=3)
        return tagFrame

    def tag_btn(self, parent, cmd):
        tagBtn = ttk.Button(parent, text='Tag', command=cmd)
        row = self.lastTagRow + 1
        col = 0
        colSpan = self.lastTagCol + 1
        tagBtn.grid(row=row, column=col, columnspan=colSpan)
        return tagBtn

    def tag_checkboxes(self, parent):
        tagCheckBoxes = []
        for i, opt in enumerate(self.cfg.options('tags')):
            self.lastTagRow = i
            tags = self.cfg.get('tags', opt)
            tags = tags.split(',')
            for j, tag in enumerate(tags):
                val = tk.StringVar(value='0')
                cbtn = ttk.Checkbutton(parent, text=tag, variable=val)
                cbtn.meta = {'tag': tag, 'val': val}
                tagCheckBoxes.append(cbtn)
                cbtn.grid(row=i, column=j, sticky=tk.W)
                if self.lastTagCol < j:
                    self.lastTagCol = j
        return tagCheckBoxes

    def load_data(self):
        data = pd.read_csv(self.dataCsvPath, index_col='img') if path.exists(
            self.dataCsvPath) else pd.DataFrame()
        return data

    def update_data_from_file(self):
        self.dataFrame = self.load_data()

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

    def update_tags(self):
        imgName = self.imgGenerator.get_img_name()
        try:
            data = self.dataFrame.loc[imgName]
            for cbtn in self.tagCheckBoxes:
                try:
                    v = data[cbtn.meta['tag']]
                    v = int(data[cbtn.meta['tag']]) if not pd.isna(v) else 0
                    cbtn.meta['val'].set(v)
                except KeyError as e:
                    continue
        except KeyError:
            for cbtn in self.tagCheckBoxes:
                cbtn.meta['val'].set(0)

    def update_all(self):
        self.update_img()
        self.update_index()
        self.update_path()
        self.update_tags()
        self.update_score_label()

    def file_open_action(self):
        newDir = filedialog.askdirectory(initialdir=self.workingDir,
                                         mustexist=True)
        if newDir == '':
            return
        self.workingDir = newDir
        self.dataCsvPath = path.join(self.workingDir, 'data.csv')
        self.imgGenerator.set_working_dir(self.workingDir)
        self.currImg = self.imgGenerator.get_img()
        self.imgGenerator.set_index(0)
        self.update_all()

    def next_action(self):
        self.currImg = self.imgGenerator.next()
        self.update_all()

    def prev_action(self):
        self.currImg = self.imgGenerator.prev()
        self.update_all()

    def dedupe_action(self):
        self.imgGenerator.dedupe()
        self.update_all()

    def save_csv_action(self):
        self.dataFrame.to_csv(self.dataCsvPath, index_label='img')

    def merge_action(self):
        newDir = filedialog.askdirectory(initialdir=self.workingDir,
                                         mustexist=True)
        if newDir == '':
            return
        newDataFrame = pd.read_csv(path.join(newDir, 'data.csv'),
                                   index_col='img')
        self.dataFrame = pd.concat([self.dataFrame, newDataFrame])
        self.dataFrame = self.dataFrame[~self.dataFrame.index.duplicated(
            keep='last')]
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

    def reset_action(self):
        os.remove(self.dataCsvPath)
        self.update_data_from_file()
        self.imgGenerator.set_working_dir(self.workingDir)
        self.imgGenerator.set_index(0)
        self.update_all()

    def tag_action(self):
        for cb in self.tagCheckBoxes:
            tag = cb.meta['tag']
            val = cb.meta['val'].get()
            if val == '1':
                self.dataFrame.loc[self.imgGenerator.get_img_name(),
                                   tag] = float(val)
            else:
                self.dataFrame.loc[self.imgGenerator.get_img_name(),
                                   tag] = np.nan
        self.next_untaged()
        self.save_csv_action()

    def next_untaged(self):
        for i, imgName in enumerate(self.imgGenerator.imgNameList):
            try:
                self.dataFrame.loc[imgName]
            except KeyError:
                self.imgGenerator.set_index(i)
                self.update_all()
                return

    ######################## score ##############################
    def load_score_data(self):
        data = pd.read_csv(self.scoreCsvPath, index_col='img') if path.exists(
            self.scoreCsvPath) else pd.DataFrame(
                columns=['img', 'score']).set_index('img')
        return data

    def set_score(self, score):
        imgName = self.imgGenerator.get_img_name()
        self.scoreData.loc[imgName, 'score'] = score
        self.scoreData.to_csv(self.scoreCsvPath, index_label='img')
        self.next_unscored()

    def next_unscored(self):
        for i, imgName in enumerate(self.imgGenerator.imgNameList):
            try:
                self.scoreData.loc[imgName]
            except KeyError:
                self.imgGenerator.set_index(i)
                self.update_all()
                return

    def score_label(self, parent):
        scoreLabel = ttk.Label(parent, text='score: 0')
        scoreLabel.pack(side=tk.LEFT)
        return scoreLabel

    def update_score_label(self):
        try:
            score = self.scoreData.loc[self.imgGenerator.get_img_name(),
                                       'score']
            self.scoreLabel.configure(text=f'score: {score}')
        except KeyError:
            self.scoreLabel.configure(text='score: 0')

    ######################## score ##############################


app = App()
app.mainloop()
