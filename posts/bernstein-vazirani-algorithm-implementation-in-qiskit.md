# Bernstein-Vazirani algorithm implementation in Qiskit
It is a restricted version of the Deutsch-Jozsa algorithm where instead of distinguishing between two different classes of functions, it tries to learn a string encoded in a function.

Here is implementation for Bernstein-Vaziarini algorithm:
```python
from qiskit import *
%matplotlib inline
from qiskit.tools.visualization import plot_histogram
```


```python
num = '101001'
circuit = QuantumCircuit(len(num)+1, len(num))
circuit.h(range(len(num)))
circuit.x(len(num))
circuit.h(len(num))
circuit.barrier()
```




    <qiskit.circuit.instructionset.InstructionSet at 0x6370a60b4e80>




```python
for i, val in enumerate(reversed(num)):
    if val == "1":
        circuit.cx(i,6)
```


```python
circuit.barrier()
circuit.h(range(len(num)))
```




    <qiskit.circuit.instructionset.InstructionSet at 0x6370a60b4e20>




```python
circuit.measure(range(len(num)), range(len(num)))
```




    <qiskit.circuit.instructionset.InstructionSet at 0x6370e006d5a0>




```python
circuit.draw(output='mpl')
```




    
![png](/images/output_5_0.png)
    




```python
sim = Aer.get_backend('qasm_simulator')
result = execute(circuit, sim, shots = 1).result()
counts = result.get_counts()
print(counts)
```

    {'101001': 1}


