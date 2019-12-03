# Bribe the Prisoners
In a kingdom there are prison cells (numbered $1$ to $P$) built to form a straight line segment.
Cells number $i$ and $i+1$ are adjacent, and prisoners in adjacent cells are called "neighbours."
A wall with a window separates adjacent cells, and neighbours can communicate through that window.

All prisoners live in peace until a prisoner is released.
When that happens, the released prisoner's neighbours find out, and each communicates this to his other neighbour.
That prisoner passes it on to his other neighbour, and so on until they reach a prisoner with no other neighbour (because he is in cell $1$, or in cell $P$, or the other adjacent cell is empty).
A prisoner who discovers that another prisoner has been released will angrily break everything in his cell, unless he is bribed with a gold coin.
So, after releasing a prisoner in cell $A$, all prisoners housed on either side of cell $A$ - until cell $1$, cell $P$ or an empty cell -- need to be bribed.

Assume that each prison cell is initially occupied by exactly one prisoner, and that only one prisoner can be released per day.
Given the list of $Q$ prisoners to be released in $Q$ days, find the minimum total number of gold coins needed as bribes if the prisoners may be released in any order.

Note that each bribe only has an effect for one day.
If a prisoner who was bribed yesterday hears about another released prisoner today, then he needs to be bribed again.

## Input
The first line of the input contains the two integers $P$ ($1 \le P \le 10\,000$) and $Q$ ($1 \le Q \le \min(P, 100)$).

The next line contains distinct $Q$ integers separated by spaces -- the cell numbers of the prisoners to be released.
The cell numbers are between $1$ and $P$, and sorted in ascending order.

## Output
Output a single integer, the minimum number of gold coins needed as bribes.

## Scoring
Your solution will be tested on a set of test groups, each worth a number of points.
To get the points for a test group you need to solve all test cases in the test group. Your final score will be the maximum score of a single submission.

| Group | Points | Constraints|
| :---- |:-------| :----------|
| 1     | 15     | $P \le 100$, $Q \le 5$ |
| 2     | 35     | No additional constraints. |
