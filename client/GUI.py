from Tkinter import Tk, Text, BOTH, W, N, E, S
from ttk import Frame, Button, Label, Style
import subprocess


class Example(Frame):
  
    def __init__(self, parent):
        Frame.__init__(self, parent)   
         
        self.parent = parent
        
        self.initUI()
        
    def initUI(self):
      
        self.parent.title("Windows")
        self.style = Style()
        self.style.theme_use("default")
        self.pack(fill=BOTH, expand=1)

        self.columnconfigure(1, weight=1)
        self.columnconfigure(3, pad=7)
        self.rowconfigure(3, weight=1)
        self.rowconfigure(5, pad=7)
        
        lbl = Label(self, text="Windows")
        lbl.grid(sticky=W, pady=4, padx=5)
        
        area = Text(self)
        area.grid(row=1, column=0, columnspan=2, rowspan=4, 
            padx=5, sticky=E+W+S+N)
        
        abtn = Button(self, text="Push", command=subprocess.call(#push))
        abtn.grid(row=1, column=3)

        cbtn = Button(self, text="Pull", command=subprocess.call(#pull))
        cbtn.grid(row=2, column=3, pady=4)
        
        hbtn = Button(self, text="Close", command=self.quit)
        hbtn.grid(row=3, column=0, padx=5)         

def main():
  
    root = Tk()
    root.geometry("350x300+300+300")
    app = Example(root)
    root.mainloop()  


if __name__ == '__main__':
    main() 