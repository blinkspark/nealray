import os, sys, re, shutil
import tkinter as tk
from tkinter import ttk, filedialog
from PIL import Image, ImageTk
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
        ttk.Style().configure('My.TCheckbutton', font=(18))

        self.tagNames = {
            'R18': tk.StringVar(value='0'),
            '2D': tk.StringVar(value='0'),
            '3D': tk.StringVar(value='0'),
            # '现实': tk.StringVar(value='0'),
            '差': tk.StringVar(value='0'),
            '中': tk.StringVar(value='0'),
            '好': tk.StringVar(value='0'),
            'BDSM': tk.StringVar(value='0'),
            'Latex/Rubber': tk.StringVar(value='0'),
            'Messy': tk.StringVar(value='0'),
            'Furry': tk.StringVar(value='0'),
            '兽耳': tk.StringVar(value='0'),
            '变装': tk.StringVar(value='0'),
            '裸露': tk.StringVar(value='0'),
            '性器官': tk.StringVar(value='0'),
            '性交': tk.StringVar(value='0'),
            '性玩具': tk.StringVar(value='0'),
            '触手': tk.StringVar(value='0'),
            '暴露': tk.StringVar(value='0'),
            'ABDL': tk.StringVar(value='0'),
        }

        self.workingDir = sys.argv[1] if (
            len(sys.argv) > 1 and os.path.exists(sys.argv[1])) else '.'

        self.dataCsvPath = path.join(self.workingDir, 'data.csv')
        self.dataFrame = self.load_data()

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
        self.tagFrame = self.tag_frame(self.mainFrame)
        self.tagCheckBoxes = self.tag_checkboxes(self.tagFrame)
        self.tagBtn = self.tag_btn(self.tagFrame, self.tag_action)

        self.next_untaged()
        # self.update_all()

        self.bind('<Right>', lambda e: self.next_action())
        self.bind('<Left>', lambda e: self.prev_action())
        self.bind(
            '<Home>', lambda e:
            (self.imgGenerator.set_index(0), self.update_all()))
        self.bind('<Escape>', lambda e: self.quit())

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

    @staticmethod
    def tag_btn(parent, cmd):
        tagBtn = ttk.Button(parent, text='Tag', command=cmd)
        tagBtn.pack(side=tk.BOTTOM)
        return tagBtn

    def tag_checkboxes(self, parent):
        tagCheckBoxes = []
        for tagName in self.tagNames:
            # font size 20
            tagCheckbox = ttk.Checkbutton(parent,
                                          padding=(5, 5),
                                          text=tagName,
                                          style='My.TCheckbutton',
                                          variable=self.tagNames[tagName])
            tagCheckbox.pack(side=tk.LEFT)
            tagCheckBoxes.append(tagCheckbox)
        return tagCheckBoxes

    def load_data(self):
        data = pd.read_csv(self.dataCsvPath, index_col='id') if path.exists(
            self.dataCsvPath) else pd.DataFrame(
                columns=['img'].extend(list(self.tagNames)))
        if len(data.columns) > 0:
            for tagName in self.tagNames:
                if not tagName in data.columns:
                    data[tagName] = pd.Series([], dtype=pd.Int64Dtype)
                    data[tagName].fillna(0, inplace=True)
        for c in data.columns:
            print(data[c].dtype)
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
        if self.dataFrame.eq(imgName).any().any():
            data = self.dataFrame.loc[self.dataFrame.eq(imgName).any(axis=1)]
            for tagName in self.tagNames:
                self.tagNames[tagName].set(data[tagName].values[0])
        else:
            for tagName in self.tagNames:
                self.tagNames[tagName].set('0')

    def update_all(self):
        self.update_img()
        self.update_index()
        self.update_path()
        self.update_tags()

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

    def dedupe_data(self):
        self.dataFrame.drop_duplicates(subset='img', inplace=True, keep='last')
        self.dataFrame.reset_index(drop=True, inplace=True)

    def save_csv_action(self):
        self.dataFrame.to_csv(self.dataCsvPath, index_label='id')

    def merge_action(self):
        newDir = filedialog.askdirectory(initialdir=self.workingDir,
                                         mustexist=True)
        if newDir == '':
            return
        newDataFrame = pd.read_csv(path.join(newDir, 'data.csv'),
                                   index_col='id')
        self.dataFrame = pd.concat([self.dataFrame, newDataFrame],
                                   ignore_index=True)
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

    def reset_action(self):
        os.remove(self.dataCsvPath)
        self.update_data_from_file()
        self.imgGenerator.set_working_dir(self.workingDir)
        self.imgGenerator.set_index(0)
        self.update_all()

    def tag_action(self):
        tagDict = {
            'img': self.imgGenerator.imgNameList[self.imgGenerator.index]
        }
        for tagName in self.tagNames:
            tagDict[tagName] = int(self.tagNames[tagName].get())
        self.dataFrame = pd.concat([
            self.dataFrame,
            pd.DataFrame(tagDict, index=[self.imgGenerator.index])
        ],
                                   ignore_index=True)
        self.dedupe_data()
        self.next_untaged()

    def next_untaged(self):
        for i, imgName in enumerate(self.imgGenerator.imgNameList):
            try:
                if not self.dataFrame['img'].eq(imgName).any():
                    self.imgGenerator.set_index(i)
                    self.update_all()
                    return
            except Exception as e:
                # self.imgGenerator.next()
                self.update_all()
                return


app = App()
app.mainloop()
