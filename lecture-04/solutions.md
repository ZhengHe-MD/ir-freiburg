# Solutions

## Exercise 1

$$
\begin{align*}
\mathcal{L} &= \sum_{i=1}^{m}p_{i}L_{i} + \lambda(\sum_{i=1}^{m}2^{-L_{i}} - 1) \\
\frac{\partial \mathcal{L}}{\partial \lambda} &= \sum_{i=1}^{m}2^{-L{i}} - 1 = 0 \Rightarrow \lambda = \frac{1}{ln2} \\
\frac{\partial \mathcal{L}}{\partial L_{i}} &= p_{i} + \lambda(-2^{-L{i}}ln2) = 0 \Rightarrow L_i = -log_{2}{p_{i}}\\
&\Rightarrow min(E(L_{x})) = \sum_{i=1}^{m}-p_{i}log_2{p_i} = H(X)
\end{align*}
$$

But I don't know how to prove  it is a minimum. It remains to be worked out.

## Exercise 2

For a Golomb encoding, \#bits taken for encoding each gap is
$$
\begin{align*}
L_{i} &= \left \lfloor \frac{ip}{ln2} \right \rfloor + 1 + \left \lceil log_2{\frac{ln2}{p}} \right \rceil \\
&\le \frac{ip}{ln2} + 1 + log_2{(\frac{ln2}{p} + 1)} \\
&\le \frac{ip}{ln2} + 1 + log_2{\frac{3ln2}{p}} \\
&=\frac{ip}{ln2} + 1 + log_2{\frac{1}{p}} + log_2{3ln2}
\end{align*}
$$
We want to prove that
$$
\begin{align*}
L_i &\le log_{2}{\frac{1}{p_i}} + O(1) \\
&= log_{2}{\frac{1}{(1-p)^{i-1}p}} + O(1) \\
&= log_{2}(\frac{1}{1-p})^{i-1} + log_{2}(\frac{1}{p}) + O(1) \\
\end{align*}
$$
Compare both sides, we get
$$
\begin{align*}
\frac{ip}{ln2} + 1 + log_2{\frac{1}{p}} + log_2{3ln2} &\le log_{2}(\frac{1}{1-p})^{i-1} + log_{2}(\frac{1}{p}) + O(1) \\
\frac{ip}{ln2} + 1 + log_2{3ln2} &\le log_{2}(\frac{1}{1-p})^{i-1} + O(1) \\
\frac{(i-1)p}{ln2} + \frac{p}{ln2} + 1 + log_2{3ln2} &\le (i-1)log_2{\frac{1}{1-p}} + O(1) \\
\end{align*}
$$
Since $\frac{p}{ln2} + 1 + log_2{3ln2} < 6$, we only need to show
$$
\begin{align*}
\frac{p}{ln2} &\le log_2{\frac{1}{1-p}} \\
p &\le ln\frac{1}{1-p} \\
e^p &\le \frac{1}{1-p} \\
e^{-p} &\ge 1-p 
\end{align*}
$$
which is given by the hint. Q.E.D

## Exercise 3

We assume that $|L_j|$ is propotional to $\frac{1}{j}$, say $|L_j| = \frac{k}{j}$,
$$
\begin{align*}
\sum_{j=1}^{m}|L_i| &= N \\
k\sum_{j=1}^{m}\frac{1}{j} &= N \\
&\Rightarrow k = \frac{N}{lnm + O(1)}
\end{align*}
$$
The expected total number of bits required to gap-encode all the inverted lists is
$$
\begin{align*}
E &= \sum_{j=1}^{m}(log_2j + O(1))|L_i|\\
&= \sum_{j=1}^{m}|Li|log_2j + O(N)
\end{align*}
$$
We want to prove that
$$
\begin{align*}
\sum_{j=1}^{m}|Li|log_2j + O(N) &\le N\frac{log_2m}{2} + O(N)
\end{align*}
$$
Since $|L_j| = \frac{k}{j}$, the above inequality becomes
$$
\begin{align*}
\frac{N}{ln2(lnm+O(1))}\sum_{i=1}^{m}\frac{lnj}{j} + O(N) \le N\frac{log_2m}{2} + O(N)\\
\frac{N}{ln2(lnm+O(1))}(\frac{ln^2m}{2}+O(1)) + O(N) \le N\frac{log_2m}{2} + O(N)\\
\end{align*}
$$
Start from the left
$$
\begin{align*}
\frac{N}{ln2(lnm+O(1))}(\frac{ln^2m}{2}+O(1)) + O(N) &\le \frac{N}{(ln2)(lnm)}(\frac{ln^2m}{2}+O(1)) + O(N) \\
&\le N\frac{log_2m}{2} + \frac{N}{(ln2)(lnm)}O(1) + O(N) \\
&\le N\frac{log_2m}{2} + O(N)
\end{align*}
$$
Q.E.D