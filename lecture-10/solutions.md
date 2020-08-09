# Solutions

## Exercise 1

$$
\begin{align*}
U \cdot S \cdot V &= A \\
U^T \cdot U \cdot S \cdot V &= U^T \cdot A \\
S \cdot V &= U^T \cdot A \\
V &= S^{-1} \cdot U^T \cdot A
\end{align*}
$$

## Exercise 2

##### Part 1

$$
\begin{align*}
A \cdot A^T = 
\begin{pmatrix}
13&5&5&0&13 \\
9&15&15&0&9 \\
0&0&0&20&0 \\
\end{pmatrix} \cdot
\begin{pmatrix}
13&9&0 \\
5&15&0 \\
5&15&0 \\
0&0&20 \\
13&9&0 \\
\end{pmatrix}
\end{align*} = 
\begin{pmatrix}
388 & 384 & 0 \\
384 & 612 & 0 \\
0 & 0 & 400
\end{pmatrix}
$$

##### Part 2

1. Guess the first Eigenvector and the corresponding Eigenvalue

$$
\begin{align*}
\begin{pmatrix}
388 & 384 & 0 \\
384 & 612 & 0 \\
0 & 0 & 400
\end{pmatrix} \cdot
\begin{pmatrix}
0 \\
0 \\
1
\end{pmatrix} = 400 \cdot
\begin{pmatrix}
0 \\
0 \\
1
\end{pmatrix}
\end{align*}
$$

2. Compute the rest Eigenvalues

$$
\begin{align*}
\begin{vmatrix}
388 - \lambda & 384 \\
384 & 612 - \lambda 
\end{vmatrix} &= 0 \\
\lambda^2 - 1000\lambda + 90000 &= 0 \\
\lambda_1 = 900, \lambda_2 = 100
\end{align*}
$$

3. Figure out the corresponding Eigenvectors

$$
\begin{align*}
\begin{pmatrix}
388 & 384 \\
384 & 612
\end{pmatrix} \cdot
\begin{pmatrix}
-\frac{4}{5} \\
\frac{3}{5}
\end{pmatrix} = 100 \cdot
\begin{pmatrix}
-\frac{4}{5} \\
\frac{3}{5}
\end{pmatrix}
\end{align*}
$$

$$
\begin{align*}
\begin{pmatrix}
388 & 384 \\
384 & 612
\end{pmatrix} \cdot
\begin{pmatrix}
\frac{3}{5} \\
\frac{4}{5}
\end{pmatrix} = 900 \cdot
\begin{pmatrix}
\frac{3}{5} \\
\frac{4}{5}
\end{pmatrix}
\end{align*}
$$

The EVD of $A \cdot A^T$ is
$$
\begin{align*}
A \cdot A^T =
\begin{pmatrix}
388 & 384 & 0 \\
384 & 612 & 0 \\
0 & 0 & 400
\end{pmatrix} = 
\begin{pmatrix}
\frac{3}{5} & 0 & -\frac{4}{5} \\
\frac{4}{5} & 0 & \frac{3}{5} \\
0 & 1 & 0
\end{pmatrix} \cdot
\begin{pmatrix}
900 & 0 & 0 \\
0 & 400 & 0 \\
0 & 0 & 100
\end{pmatrix} \cdot
\begin{pmatrix}
\frac{3}{5} & \frac{4}{5} & 0 \\
0 & 0 & 1 \\
-\frac{4}{5} & \frac{3}{5} & 0
\end{pmatrix}
\end{align*}
$$

##### Part 3

$$
U =
\begin{pmatrix}
\frac{3}{5} & 0 & -\frac{4}{5} \\
\frac{4}{5} & 0 & \frac{3}{5} \\
0 & 1 & 0
\end{pmatrix},
U \cdot U^T = 
\begin{pmatrix}
1 & 0 & 0 \\
0 & 1 & 0 \\
0 & 0 & 1
\end{pmatrix}
$$

$U$ is indeed a $3 \times 3$ column-orthonormal matrix.
$$
S = 
\begin{pmatrix}
30 & 0 & 0 \\
0 & 20 & 0 \\
0 & 0 & 10
\end{pmatrix}
$$

##### Part 4

$$
\begin{align*}
V &= S^{-1}\cdot U^T \cdot A \\
&= 
\begin{pmatrix}
\frac{1}{30} & 0 & 0 \\
0 & \frac{1}{20} & 0 \\
0 & 0 & \frac{1}{10}
\end{pmatrix} \cdot
\begin{pmatrix}
\frac{3}{5} & \frac{4}{5} & 0 \\
0 & 0 & 1 \\
-\frac{4}{5} & \frac{3}{5} & 0
\end{pmatrix} \cdot
\begin{pmatrix}
13&5&5&0&13 \\
9&15&15&0&9 \\
0&0&0&20&0 \\
\end{pmatrix}
=
\begin{pmatrix}
\frac{1}{2} & \frac{1}{2} & \frac{1}{2} & 0 & \frac{1}{2} \\
0 & 0 & 0 & 1 & 0 \\
-\frac{1}{2} & \frac{1}{2} & \frac{1}{2} & 0 & -\frac{1}{2}
\end{pmatrix}
\end{align*}
$$

##### Part 5

$$
U \cdot S \cdot V =
\begin{pmatrix}
\frac{3}{5} & 0 & -\frac{4}{5} \\
\frac{4}{5} & 0 & \frac{3}{5} \\
0 & 1 & 0
\end{pmatrix} \cdot
\begin{pmatrix}
30 & 0 & 0 \\
0 & 20 & 0 \\
0 & 0 & 10
\end{pmatrix} \cdot
\begin{pmatrix}
\frac{1}{2} & \frac{1}{2} & \frac{1}{2} & 0 & \frac{1}{2} \\
0 & 0 & 0 & 1 & 0 \\
-\frac{1}{2} & \frac{1}{2} & \frac{1}{2} & 0 & -\frac{1}{2}
\end{pmatrix} =
\begin{pmatrix}
13&5&5&0&13 \\
9&15&15&0&9 \\
0&0&0&20&0 \\
\end{pmatrix} = A
$$

## Exercise 3

##### Variant 1

$$
A_k = U_k \cdot S_k \cdot V_k =
\begin{pmatrix}
\frac{3}{5} & 0 \\
\frac{4}{5} & 0 \\
0 & 1
\end{pmatrix} \cdot
\begin{pmatrix}
30 & 0 \\
0 & 20
\end{pmatrix} \cdot
\begin{pmatrix}
\frac{1}{2} & \frac{1}{2} & \frac{1}{2} & 0 & \frac{1}{2} \\
0 & 0 & 0 & 1 & 0
\end{pmatrix} =
\begin{pmatrix}
9 & 9 & 9 & 0 & 9 \\
12 & 12 & 12 & 0 & 12 \\
0 & 0 & 0 & 20 & 0
\end{pmatrix}
$$


$$
q^T \cdot A_k =
\begin{pmatrix}
4 & 1 & 2
\end{pmatrix} \cdot
\begin{pmatrix}
9 & 9 & 9 & 0 & 9 \\
12 & 12 & 12 & 0 & 12 \\
0 & 0 & 0 & 20 & 0
\end{pmatrix} = 
\begin{pmatrix}
48 & 48 & 48 & 40 & 48
\end{pmatrix}
$$

##### Variant 2

$$
q_k^{T} = q^T \cdot U_k \cdot S_k = 
\begin{pmatrix}
4 & 1 & 2
\end{pmatrix} \cdot
\begin{pmatrix}
\frac{3}{5} & 0 \\
\frac{4}{5} & 0 \\
0 & 1
\end{pmatrix} \cdot
\begin{pmatrix}
30 & 0 \\
0 & 20
\end{pmatrix} = 
\begin{pmatrix}
96 & 40
\end{pmatrix}
$$

$$
q_k^T \cdot V_k =
\begin{pmatrix}
96 & 40
\end{pmatrix} \cdot
\begin{pmatrix}
\frac{1}{2} & \frac{1}{2} & \frac{1}{2} & 0 & \frac{1}{2} \\
0 & 0 & 0 & 1 & 0
\end{pmatrix} =
\begin{pmatrix}
48 & 48 & 48 & 40 & 48
\end{pmatrix}
$$

##### Variant 3

$$
T_k = U_k \cdot U_k^{T} =
\begin{pmatrix}
\frac{3}{5} & 0 \\
\frac{4}{5} & 0 \\
0 & 1
\end{pmatrix} \cdot
\begin{pmatrix}
\frac{3}{5} & \frac{4}{5} & 0 \\
0 & 0 & 1
\end{pmatrix} = 
\begin{pmatrix}
\frac{9}{25} & \frac{12}{25} & 0 \\
\frac{12}{25} & \frac{16}{25} & 0 \\
0 & 0 & 1
\end{pmatrix}
$$

use $\frac{12}{25}$ as the 0-1 threshold：
$$
T_k^{'} = 
\begin{pmatrix}
0 & 1 & 0 \\
1 & 1 & 0 \\
0 & 0 & 1
\end{pmatrix}
$$

$$
q^T \cdot T_k^{'} \cdot A =
\begin{pmatrix}
4 & 1 & 2
\end{pmatrix} \cdot
\begin{pmatrix}
0 & 1 & 0 \\
1 & 1 & 0 \\
0 & 0 & 1
\end{pmatrix} \cdot
\begin{pmatrix}
13&5&5&0&13 \\
9&15&15&0&9 \\
0&0&0&20&0 \\
\end{pmatrix} = 
\begin{pmatrix}
58 & 80 & 80 & 40 & 58
\end{pmatrix}
$$

