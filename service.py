from multiprocessing import Process, Pipe
import subprocess


def discovery_function(db_out):
    discovery_shell = subprocess.Popen(["./test.sh"], stdout=subprocess.PIPE)
    db_out.send(stdout.communicate())
    db_out.close()

if __name__=='__main__':
    out_changes, in_changes = Pipe()
    discovery= Process(target=discovery_function, args=(in_changes,))
    discovery.start()
    print out_changes.recv()
    discovery.join()
