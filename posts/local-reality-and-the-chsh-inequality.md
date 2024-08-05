# Local Reality and the CHSH inequality
The CHSH game is a thought experiment involving two parties separated at a great distance (far enough to preclude classical communication at the speed of light), each of whom has access to one half of an entangled two-qubit pair. Analysis of this game shows that no classical local hidden-variable theory can explain the correlations that can result from entanglement. Since this game is indeed physically realizable, this gives strong evidence that classical physics is fundamentally incapable of explaining certain quantum phenomena, at least in a "local" fashion. 

Here is implementation for CHSH Inequality
```python
    import qiskit
    from qiskit import QuantumCircuit, ClassicalRegister, QuantumRegister,transpile, Aer
    from qiskit.tools.visualization import circuit_drawer
    from qiskit.tools.monitor import job_monitor, backend_monitor, backend_overview
    from qiskit.providers.aer import noise
    from qiskit_ibm_provider import IBMProvider

    import matplotlib.pyplot as plt
    import numpy as np
    import time
```


```python
    sim = Aer.get_backend('aer_simulator')

    provider = IBMProvider('imbq-key')
    backend = provider.get_backend('ibmq_lima')
```


```python
    def make_chsh_circuit(theta_vec):
        chsh_circuits = []

        for theta in theta_vec:
            obs_vec = ['00', '01', '10', '11']
            for el in obs_vec:
                qc = QuantumCircuit(2,2)
                qc.h(0)
                qc.cx(0,1)
                qc.ry(theta, 0)
                for a in range(2):
                    if el[a] == '1':
                        qc.h(a)
                qc.measure(range(2),range(2))
                chsh_circuits.append(qc)
        return chsh_circuits
```


```python
    def compute_chsh_witness(counts):
        CHSH1 = []
        CHSH2 = []
        for i in range(0, len(counts), 4):
            theta_dict = counts[i:i + 4]
            zz = theta_dict[0]
            zx = theta_dict[1]
            xz = theta_dict[2]
            xx = theta_dict[3]

            no_shots = sum(xx[y] for y in xx)
            chsh1 = 0
            chsh2 = 0

            for element in zz:
                parity = (-1)**(int(element[0])+int(element[1]))
                chsh1+=parity*zz[element]
                chsh2+=parity*zz[element]
        
            for element in zx:
                parity = (-1)**(int(element[0])+int(element[1]))
                chsh1+= parity*zx[element]
                chsh2-= parity*zx[element]
        
            for element in xz:
                parity = (-1)**(int(element[0])+int(element[1]))
                chsh1-= parity*xz[element]
                chsh2+= parity*xz[element]
        
            for element in xx:
                parity = (-1)**(int(element[0])+int(element[1]))
                chsh1+= parity*xx[element]
                chsh2+= parity*xx[element]
        
            CHSH1.append(chsh1/no_shots)
            CHSH2.append(chsh2/no_shots)

        return CHSH1, CHSH2
```


```python
    number_of_thetas = 15
    theta_vec = np.linspace(0,2*np.pi,number_of_thetas)
    my_chsh_circuits = make_chsh_circuit(theta_vec)
```


```python
    my_chsh_circuits[4].draw(output="mpl")
```


```python
    result_ideal = sim.run(my_chsh_circuits).result()

    tic = time.time()
    transpiled_circuits = transpile(my_chsh_circuits, backend)
    job_real = backend.run(transpiled_circuits, shots=8192)
    job_monitor(job_real)
    result_real = job_real.result()
    toc = time.time()
    print(toc-tic)
```


```python
    CHSH1_ideal, CHSH2_ideal = compute_chsh_witness(result_ideal.get_counts())
    CHSH1_real, CHSH2_real = compute_chsh_witness(result_real.get_counts())
```


```python
    plt.figure(figsize=(12,8))
    plt.plot(theta_vec, CHSH1_ideal,'o-',label = 'CHSH1 Noiseless')
    plt.plot(theta_vec, CHSH2_ideal, 'o-', label='CHSH2 Noiseless')

    plt.plot(theta_vec, CHSH1_real, 'x-', label='CHSH1 Quito')
    plt.plot(theta_vec, CHSH2_real, 'x-', label='CHSH2 Quito')

    plt.legend()
```


