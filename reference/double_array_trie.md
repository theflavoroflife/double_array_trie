# 双数字典树(double array trie)



## 前言

这篇文章主要是介绍一种屏蔽字快速检测的方法：双数组字典树。

主要内容包括：

- 双数组字典树的base/check数组构造流程，额外维护了一个lengths数组。
- fail数组的建立流程。
- 给出一个示例，演示一遍双数组字典树的查询流程。

双数组字典树是**多模式匹配**的一种实现方式，利用base和check数组做状态转移，理论上的查找速度比基于哈希映射的trie的查找速度更快。



## 双数组字典树(double array trie)构建base和check数组

问题：给定多个模式串：i，he，his，she，hers，构建这些模式串的双数组字典树的base和check数组。

##### 一、构建流程前言

- 构建base和check的一个小技巧：对给定的多个模式串进行**字典序的升序排序**，好处是可以通过下标找到孩子节点，不需要使用[论文](https://kampersanda.github.io/pdf/KAIS2017.pdf)中的tail数组做辅助。

- 本次演示增加一个lengths数组记录结束字符的时候模式串的长度，这样做的好处是可以抛弃[论文](https://kampersanda.github.io/pdf/KAIS2017.pdf)中使用负数表示模式串结束的方法，同时也能记录多个以当前字符为结束字符时模式串的长度。

##### 二、构建base和check的规则

- 构建前的操作是，将给定的多个模式串[i，he，his，she，hers]进行从小到大排序，排序后：

![avatar](./images/sorted_words.jpg)

- 构建base和check数组的时候遍历这些模式串的规则是：层次遍历(BFS)。

- 理解下面2个公式的含义很重要，构造base和check数组的时候会不断的使用到这2个规则。

```
  公式一：base[s] + c = t
  s：      理解为父亲节点的坐标位置(index)
  base[s]：表示位置s对应父节点存储的base值，base[s]不为空说明这个节点有孩子节点。
  c：      对应孩子节点的ascii码
  t：      计算出来孩子节点应该放入的位置(index)，通过base[t]或者check[t]检查是否可以放入这个位置。
  如果check[t]的结果为空，表示寻找到的孩子节点的位置(index)未被占用，此孩子节点c可以暂时放在t这个位置。
  如果check[t]的结果不为空，表示根据父节点s找到c的位置被其他节点占用，需要增加父节点base[s]的值重新寻找一个t为空的值。
  (程序实现的时候，是查看base/check/lengths数组中只要其中一个有值，就表示该位置已经被占用)
```

``` 
  公式二：check[t] = s
  t： 孩子节点的index
  s： 父节点的index
  建立孩子节点和父节点对应关系，同时也说明check[t]不为空也表示该位置t已经被占用。
  (程序实现的时候，是查看base/check/lengths数组中只要其中一个有值，就表示该位置已经被占用)
```

##### 三、构建流程

除了构造base和check数组外，一起填写的数组如下(fail数组暂时不需要)：

```
index:   数组下标，辅助理解
base:    可以理解为存储父亲节点的base值；base不为空，说明该节点有孩子节点
check:   主要建立孩子节点和父亲节点的对应关系，check存储的值就是父亲节点的位置；同时也可以用来表示位置是否被占用
fail:    查询失败后的跳转
lengths: 如果这个位置是结束字符，用来记录结束字符的长度，可以记录多个结束字符长度
```

###### 1、初始设置root节点位于index[0]处，并且初始化base[0]=1

![avatar](./images/bc_root.jpg)

###### 2、第1层：root下有孩子节点h、i、s，对应的ascii码值是h:104，i:105，s:115

![](./images/bc_first_code.jpg)

- 根据构造转移公式：base[s] + c = t，做如下操作：
  - 节点h，ascii码是104，父亲节点是root，父亲节点位于index[0]：
    - ```base[s]+c = base[0]+'h' = base[0]+104 = t(1) ```
  - 节点i，ascii码是105，父亲节点是root，父亲节点位于index[0]：
    - ```base[s]+c = base[0]+'i' = base[0]+105 = t(2)```
  - 节点s，ascii码是115，父亲节点是root，父亲节点位于index[0]：
    - ```base[s]+c = base[0]+'s' = base[0]+115 = t(3)```
  
  初始化base[0]=1，则t(1)=105，t(2)=106，t(3)=116；

  同时查看表格中check[105]、check[106]、check[116]的位置都为空，说明当base[0]=1的时候，root节点下的孩子节点h、i、s找到位置index[105]、index[106]、index[116]可以分别存放。所以**base[0]等于1**成立。
  
- 根据构造转移的公式：check[t] = s，做如下操作：
  
  - 孩子节点的check数组要指向父亲节点的下标，节点h、i、s的父亲节点都是root节点，root节点的下标是0，建立了孩子节点h、i、s和父亲节点root的对应关系：
    - **```check[105] = 0```**
    - **```check[106] = 0```**
    - **```check[116] = 0```**

- 同时这个root节点下的第二个孩子节点i是结束字符，字符i位于index[106]，需要记录**lengths[106]=1**。

![](./images/bc_first_depth.jpg)

###### 3、第2层节点：e、i、h，对应的ascii码值是e:101，i:105，h:104

![avatar](./images/bc_second_code.jpg)

- 根据构造转移的公式：base[s] + c = t，做如下操作：
  
  - 节点e、i的父亲节点是h。
    - 节点e，ascii码是101，父亲节点是h，父亲节点位置是index[105]：
      - ```base[s]+c = base[105]+'e' = base[105]+101 = t(1)```
  
    - 节点i，ascii码是105，父亲节点是h，父亲节点位置是index[105]：
      - ```base[s]+c = base[105]+'i' = base[105]+105 = t(2)```
  
    假如base[105]=1，则t(1)=102，t(2)=106。check[102]位置为空，但是check[106]位置不为空，该位置不能使用，已经被第一层节点i占用，所以base[105]=1不成立。
  
    假如base[105]=2，则t(1)=103，t(2)=107。查看表格check[103]和check[107]位置都为空，所以父节点base[105]=2的时候，孩子节点都可以找到空位置，所以**base[105]=2**。
  
  - 节点h的父亲节点是s。
    - 节点h，ascii码是104，父亲节点是s，父亲节点位置是index[116]：
      - ```base[s]+c = base[116]+'h' = base[116]+104 = t```
  
    假如base[116]=1，t=105，check[105]已经被占用；
    
    假如base[116]=2，t=106，check[106]也已经被占用；
    
    假如base[116]=3，t=107，check[107]也被同层节点i占用；
    
    假如base[116]=4，t=108，check[108]位置为空，可以使用，所以**base[116]=4**。

- 根据构造转移的公式：check[t] = s，做如下操作：
  - 孩子节点e、i的位置是103和107，需要建和父亲节点(h)的对应关系，即：
    - ```check[103]=105```
    - ```check[107]=105```
  - 孩子节点h，建立和父亲节点(s)的对应关系，即：
    - ```check[108]=116```

- 同时节点e是结束字符，字符e位于index[103]，需要记录**lengths[103]=2**。

![avatar](./images/bc_second_depth.jpg)

###### 4、第3层节点：r、s、e，对应的ascii码的值是r:114，s:115，e:101，其中e下有孩子节点r，i下有孩子节点s，h下有孩子节点e

![avatar](./images/bc_third_code.jpg)

- 根据转移的公式：base[s] + c = t，做如下操作：
  - 节点r，ascii码是114，父亲节点是e，父亲节点坐标的是index[103]
    - ```base[s]+c = base[103]+'r' = base[103]+114 = t```

    如果父亲节点e的base[103]=1，则t=115，check[115]目前位置为空，所以孩子节点r可以放在index[115]位置，最终**base[103]=1**

  - 节点s，ascii码是115，父亲节点是i，父亲节点坐标是index[107]
    - ```base[s]+c = base[107]+'s' = base[107]+115 = t```

    如果父亲节点i的base[107]=1，则t=116，check[116]目前位置不为空，所以孩子节点s不可以放在index[116]位置；
    如果父亲节点i的base[107]=2，则t=117，check[117]目前位置为空，所以孩子节点s可以放在index[117]位置，最终**base[107]=2**。

  - 节点e，ascii码是101，父亲节点是h，父亲节点的坐标是index[108]
    - ```base[s]+c = base[108]+'e' = base[108]+101 = t```

    如果父亲节点h的base[108]=1，则t=102，check[102]目前位置为空，所以孩子节点e可以放在index[102]位置，最终**base[108]=1**。

- 根据构造的公式check[t] = s：，做如下操作：
  - base[103]=1，父亲节点e找到孩子节点r可以放入的空闲位置是index[115]，建立孩子节点r和父亲节点e的对应关系，即
    - **```check[115]=103```**

  - base[107]=2，父亲节点i找到孩子节点s可以放入的空闲位置是index[117]，建立孩子节点s和父亲节点i的对应关系，即
    - **```check[117]=107```**

  - base[108]=1，父亲节点h找到孩子节点e可以放入的空闲位置是index[102]，建立孩子节点e和父亲节点h的对应关系，即
    - **```check[102]=108```**

- 同时节点s是结束字符，字符s位于index[117]，需要记录**lengths[117]=3**。

  节点e是结束字符，节点e位于index[102]，需要记录**lengths[102]=3**

![avatar](./images/bc_third_depth.jpg)

###### 5、第4层节点：s，对应的ascii码值是s:115，其中r下有孩子节点s

![avatar](./images/bc_forth_code.jpg)

- 根据构造的公式：base[s] + c = t，做如下操作：
  - 孩子节点s，ascii码是115，父亲节点是r，父亲节点位置是index[115]：
    - base[s]+c = base[115]+'s' = base[115]+115 = t
  
    如果base[115]=1，则t=116，check[116]目前位置不为空，所以孩子节点s不可以放在index[116]位置。
  
    如果base[115]=2，则t=117，check[117]目前位置不为空，所以孩子节点s不可以放在index[117]位置。
  
    如果base[115]=3，则t=118，check[118]目前位置为空，所以孩子节点s可以放在index[118]位置，最终**base[115]=3**
  
- 根据构造的公式：check[t] = s，做如下操作：
  
  - 当base[115]=3，父亲节点r，位置是index[115]，找到孩子节点s的空闲位置是index[118]，建立孩子节点s和父亲节点r的对应关系，即
    - **```check[118]=115```**

- 同时节点s是结束字符，字符s位于index[118]，需要记录**lengths[118]=4**。

![avatar](./images/bc_forth_depth.jpg)



## 双数组字典树(double array trie)构建fail数组

##### 一、构建fail数组规则：

- 按层次遍历(BFS)已经排序好的多个模式串。
- root节点和root的孩子节点(depth=1)的fail指针都指向root，也就是字符code的fail数组都填写root的index。
- 假如当前处理节点code，先找到其父节点parent，再找父节点parent的fail指针指向的节点parentFailIndex，查看节点parentFailIndex是否有与节点code相等的节点(假如child)
  - 如果找到，节点code的fail指针指向节点child，也就是节点code的fail数组填写child节点的index。
  - 如果没有找到与节点code相等的孩子节点或者parentFailIndex没有孩子节点，查看节点parentFailIndex是不是root节点
    - 如果节点parentFailIndex是root节点，则字符code的fail指针指向root。
    - 如果节点parentFailIndex不是root节点，找出节点parentFailIndex的fail指针所指向的字符，继续最开始的操作。
  
- 构造fail数组的时候，如果节点a的fail指针指向节点b，并且节点b的lengths对应的长度有值m，需要把长度m添加到节点a的lengths中。

##### 二、fail数组的构建流程

###### 1、构建root节点和第一层节点(h、i、s)的fail数组

![avatar](./images/fail_first_word.jpg)

- root节点的fail指针指向root节点，即
  - **fail[0]=0**
- 第一层节点的fail指针指向root节点，即
  - **fail[105]=0**，**fail[106]=0**，**fail[116]=0**

![avatar](./images/fail_first_depth.jpg)

###### 2、第二层节点e、i、h。 

![avatar](./images/fail_second_word.jpg)

- **节点e，index是103，ascii码是101**
  - 节点e的位置是index[**103**]，根据check[103]=105可知，节点index[103]的父亲节点是index[105]，index[105]的fail[105]=0，即父亲节点index[105]指向root节点index[0]，查表可知base[0]不为空，说明index[0]下有孩子节点，转换公式：
    - ```base[s]+c = base[0]+'e' = 1+101 = 102```
    - ```check[102] = 108```
  - 根据上面的公式可以知道，check[102] != 0，说明index[0]下没有与字符e相等的孩子节点，同时当前位置是root节点，字符e的fail指针指向root节点index[0]，即
    - **fail[103]=0**。
- **节点i，index是107，ascii码是105**
  - 节点i的位置是index[**107**]，根据check[107]=105可知，节点index[107]的父亲节点是index[105]，index[105]的fail[105]=0，即父亲节点index[105]指向root节点index[0]，查表可知base[0]不为空，说明index[0]下有孩子节点，转换公式：
  
    - ```base[s]+c = base[0]+'i' = 1+105 = 106```
  
    - ```check[106] = 0```
  - 根据上面的公式可以知道，节点index[0]下找到与字符i相等的孩子节点index[**106**]，将节点index[107]的fail指针指向节点index[106]，即
    - **fail[107]=106**
  - 查看lengths[106]=1，说明index[106]是一个结束字符，需要将lengths[106]拷贝到lengths[107]中，即
    - **lengths[107]=1**。
- **节点h，index是108，ascii码是104**
  - 节点h的位置是index[**108**]，根据check[108]=116可知，index[108]的父亲节点是index[116]，index[116]的fail[116]=0，即index[116]指向root节点index[0]，查表可知base[0]不为空，说明节点index[0]有孩子节点，转换公式：
    - ```base[s]+c = base[0]+'h' = 1+104 = 105```
    - ```check[105] = 0```
  - 根据上面的公式可以知道，节点index[0]下找到与字符h相等的孩子节点index[**105**]，将节点index[108]的fail指针指向index[105]，即
    - **fail[108]=105**。

![avatar](./images/fail_second_depth.jpg)

###### 3、第三层节点有r、s、e。

![avatar](./images/fail_third_word.jpg)

- **节点r，index是115，ascii码是114**
  - 节点r的位置是index[**115**]，根据check[115]=103可知，index[115]的父亲节点是index[103]，节点index[103]的fail[103]=0，也就是fail指针指向root节点index[0]，查表可知base[0]不为空，说明index[0]下有孩子节点，转换公式：
    - ```base[s]+c = base[0]+'r' = 1+114 = 115```
    - ```check[115] = 103```
  - 根据上面的公式可以知道，check[115] != 0，说明index[0]下没有与字符r相等的孩子节点，同时当前位置是root节点，字符e的fail指针指向root节点index[0]，即
    - **fail[115]=0**
- **节点s，index是117，ascii码是115**
  - 节点s的位置是index[**117**]，根据check[117]=107可知，index[117]的父亲节点是index[107]，节点index[107]的fail[107]=106，查表base[106]为空，说明index[106]下没有孩子节点；继续查找fail[106]=0，指向root节点index[0]，base[0]不为空，说明有孩子节点，转移公式：
      - ```base[s]+c = base[0]+'s' = 1+115 = 116```
      - ```check[116] = 0```
  - 根据上面的公式可以知道，节点index[0]下找到与字符s相等的孩子节点index[**116**]，将index[117]的fail指针指向index[116]，即
      - **fail[117]=116**
- **节点e，index是102，ascii码是101**
  - 节点e的位置是index[**102**]，根据check[102]=108可知，index[102]的父亲节点是index[108]，节点index[108]的fail[108]=105，base[105]不为空，说明index[105]下有孩子节点，
    - ```base[s]+c = base[105]+'e' = 2+101 = 103```
    - ```check[103] = 105```
  - 根据上面的公式可以知道，index[105]下找到与字符e相等的孩子节点index[**103**]，将index[102]的fail指针指向index[103]，即
    - **fail[102]=103**
  - 同时将lengths[103]的长度2添加到lengths[102]对应的lengths数组中，添加后lengths是
    - **lengths[102]=3,2**

![avatar](./images/fail_third_depth.jpg)

###### 4、第四层节点有s。

![avatar](./images/fail_forth_word.jpg)

- **节点s，index是118，ascii码是115**
  - 节点s的位置是index[**118**]，根据check[118]=115可知，节点index[118]的父亲节点是index[115]，节点index[115]的fail[115]=0，即index[115]节点的fail指针指向root节点index[0]，查表base[0]不为空，说明base[0]下有孩子节点，转移公式：
    - ```base[s]+c = base[0]+'s' = 1+115 = 116```
    - ```check[116] = 0```
  - 根据上面的公式可以知道，index[0]下找到与字符e相等的孩子节点index[**116**]，将index[118]的fail指针指向index[116]，即
    - **fail[118]=116**

![avatar](./images/fail_forth_depth.jpg)



## 双数组字典树(double array trie)查找

##### 一、查找规则

假如当前被查找字符是code(ascii码)，当前节点是current，当前的孩子节点是child。定义变量currentIndex记录当前查找位置。

- 查表查看当前节点(current)下是否有与code相等的孩子节点(child)：
  - 如果有child = code，说明找到有与字符code相等的孩子节点，则返回孩子节点(child)的下标index给currentIndex，本次查找结束。
  - 如果child != code或者当前节点没有孩子节点，则查看当前节点(current)是否是root节点：
    - 当前节点(current)是root节点，则返回root下标index给currentIndex，本次查找结束。
    - 当前节点(current)不是root节点，查表取出fail[currentIndex]的下标index给currentIndex，即currentIndex=fail[currentIndex]。也就是本次匹配失败后，转移currentIndex到fail数组指向的节点，将fail[currentIndex]位置设置为当前需要开始查找的位置，继续第一步查找流程。

再来看看下面2个公式，后面会不停的用到这2个公式的转换：

- ```base[s] + c = t```
- ```check[t] = s```

##### 二、查找流程

接下来演示一遍查找流程。

给定一段文本串"ifindhehishehersall"，查找出所有这些模式串[i，he，his，she，hers]在给定文本串中出现的位置信息。

![avatar](./images/content.jpg)

首先定义一个查找变量：**currentIndex**，记录当前在base/check的查找位置。

初次查找是从root节点开始查找，初始化currentIndex为root的index，即**currentIndex=0**。

###### 1、查找第1个字符i，ascii码是105

- 当前**currentIndex=0**。

- 查表显示base[0]不为空，说明index[0]节点下有孩子节点。

  查看index[0]节点下如果有字符i的状态转移公式，查表：

  - ```base[s]+c = base[0]+'i' = 1+105 = 106```
  - ```check[106] = 0``` 

  根据上面公式可知，index[0]节点下找到与字符i相等的孩子节点index[106]，返回孩子节点的下标106给currentIndex，即**currentIndex=106**。

- 查表显示lengths[106]=1，说明当前位置有1个结束字符，这个模式串长度是1，起始坐标是(0, 0)，匹配到的模式串是：i。

- 继续下一个字符查找。

![avatar](./images/search_first.jpg)

###### 2、查找第2个字符f，ascii码是102

- 当前**currentIndex=106**。

- 查表显示base[106]值为空，说明index[106]节点下的没有孩子节点，index[106]也不是root节点，查表显示fail[106]=0，将0赋值给currentIndex，即**currentIndex=0**。

- 查表显示base[0]不为空，表明index[0]节点下有孩子节点。

  查看index[0]节点下如果有字符f的状态转移公式，查表：

  - ```base[s]+c = base[0]+'f' = 1+102 = 103```
  - ```check[103] = 105```

  根据上面转移公式可知check[103] != 0 ，说明父亲节点index[0]下没有与字符f相等的孩子节点，返回index[0]坐标给currentIndex，即**currentIndex=0**

- 继续下一个字符查找。

![avatar](./images/search_second.jpg)

###### 3、查找第3个字符i，ascii码是105

- 当前**currentIndex=0**。

- 查表显示base[0]不为空，说明有index[0]节点下有孩子节点。

  查看index[0]节点下如果有字符i的状态转移公式，查表：

  - ```base[s]+c = base[0]+'i' = 1+105 = 106```
  - ```check[106] = 0```

  根据上面转移公式可知check[106]=0，说明父亲节点index[0]下有与字符i相等的孩子节点index[106]，返回查找到的孩子节点坐标index给currentIndex，也就是**currentIndex=106**。

- 查表显示lengths[106]=1，说明当前位置有一个结束字符，这个，模式串长度是1，起始下标是(2, 2)，匹配到的模式串是：i。

- 继续下一个字符查找。

![avatar](./images/search_third.jpg)

###### 4、查找第4个字符n，ascii码是110

- 当前**currentIndex=106**。

- 查表显示base[106]是空，说明index[106]节点下没有孩子节点，同时index[106]节点也不是root节点，查表可知fail[106]=0，将0赋值给currentIndex，即**currentIndex=0**。

- 查表可知base[0]不为空，说明index[0]节点下有孩子节点。

  查看index[0]节点下如果有字符n的状态转移公式，查表：

  - ```base[s]+c = base[0]+'n' = 1+110 = 111```
  - ```check[111] = null(空值)```

  根据上面公式可知，index[0]节点下没有与字符n相等的孩子节点，当前节点是root节点，返回root节点的index给currentIndex，即**currentIndex=0**

- 继续下一个字符查找。

![avatar](./images/search_forth.jpg)

###### 5、查找第5个字符d，ascii码是100

- 当前**currentIndex=0**。

- 查表显示base[0]不为空，说明index[0]节点下有孩子节点。

  查看index[0]节点下如果有字符d的状态转移公式，查表：

  - ```base[s]+c = base[0]+'d' = 1+100 = 101```
  - ```check[101] = null(空值)```

  根据上面的公式可以知道，index[0]节点下没有与字符d相等的孩子节点，当前index[0]节点已经是root节点，返回root的index给currentIndex，即**currentIndex=0**

- 继续下一个字符查找。

![avatar](./images/search_fifth.jpg)

###### 6、查找第6个字符h，ascii码是104

- 当前**currentIndex=0**。

- 查表显示base[0]不为空，说明index[0]节点下有孩子节点。

  查看index[0]节点下如果有字符h的状态转移公式，查表：

  - ```base[s]+c = base[0]+'h' = 1+104 = 105```
  - ```check[105] = 0```

  根据上面公式可知，index[0]节点下找到与字符h相等孩子节点index[105]，返回105给currentIndex，即**currentIndex=105**。

- 查表显示lengths[105]为空，说明当前位置没有结束字符。

- 继续下一个字符查找。

![avatar](./images/search_sixth.jpg)

###### 7、查找第7个字符e，ascii码是101

- 当前**currentIndex=105**。

- 查表显示base[105]不为空，说明index[105]节点下有孩子节点。

  查看index[105]节点下如果有字符e的状态转移公式，查表：

  - ```base[s]+c = base[105]+'e' = 2+101 = 103```
  - ```check[103] = 105```

  根据上面公式可知，index[105]节点下有与字符e相等的孩子节点index[103]，返回103给currentIndex，即**currentIndex=103**。

- 查表显示lengths[103]=2，说明当前位置有结束字符，这个模式串长度是2，起始坐标(5,6)，匹配到的模式串是：he。

- 继续下一个字符查找。

![avatar](./images/search_seventh.jpg)

###### 8、查找第8个字符h，ascii码是104

- 当前**currentIndex=103**。

- 查表显示base[103]不为空，说明index[103]节点下有孩子节点。

  查看index[103]节点下如果有字符h的状态转移公式，查表：

  - ```base[s]+c = base[103]+'h' = 1+104 = 105```
  - ```check[105] = 0```

  根据上面公式可知：check[105] != 103，说明index[103]节点下没有与字符h相等的孩子节点，同时index[103]节点也不是root节点，查表可知fail[103]=0，把0赋值给currentIndex，也就是**currentIndex=0**。

- 查表显示base[0]不为空，说明index[0]节点下有孩子节点。

  查看index[0]节点下如果有字符h的状态转移公式，查表：

  - ```base[s]+c = base[0]+'h' = 1+104 = 105```
  - ```check[105] = 0```

  根据上面的公式可知，index[0]节点找到了与h相等的孩子节点index[105]，返回105给currentIndex，即**currentIndex=105**。

- 查表显示lengths[105]为空，说明当前位置没有结束字符。

- 继续下一个字符查找。

![avatar](./images/search_eighth.jpg)

###### 9、查找第9个字符i，ascii码是105

- 当前**currentIndex=105**。

- 查表显示base[105]不为空，说明index[105]下有孩子节点。

  查看index[105]节点下如果有字符i的状态转移公式，查表：

  - ```base[s]+c = base[105]+'i' = 2+105 = 107```
  - ```check[107] = 105```

  根据上面的公式可知，index[105]下找到与字符i相等的孩子节点index[107]，返回孩子节点的index给currentIndex ，即**currentIndex =107**。

- 查表显示lengths[107]=1，说明当前位置有1个结束字符，这个字符长度为1，起始坐标(8,8)，匹配到的模式串字是：i。

- 继续下一个字符查找。

![avatar](./images/search_ninth.jpg)

###### 10、查找第10个字符s，ascii码是**115**

- 当前**currentIndex=107**。

- 查表显示base[107]不为空，说明index[107]节点下有孩子节点。

  查看index[107]节点下如果有字符s的状态转移公式，查表：

  - ```base[s]+c = base[107]+'s' = 2+115 = 117```
  - ```check[117] = 107```

  根据上面的公式可知，index[107]节点下找到与字符s相等的孩子节点index[117]，返回孩子节点的index给currentIndex ，即**currentIndex=117**。

- 查表显示lengths[117]=3，说明当前位置有1个结束字符，这个字符长度为3，起始坐标(7,9)，匹配模式字是：his。

- 继续下一个字符查找。

![avatar](./images/search_tenth.jpg)

###### 11、查找第11个字符h，ascii码是104

- 当前**currentIndex=117**。

- 查表显示base[117]为空，说明index[117]节点下没有孩子节点，查表可知fail[117]=116，将116赋值给currentIndex，即**currentIndex=116**。

- 查表显示base[116]不为空，说明index[116]节点下有孩子节点。

  查看index[116]节点下如果有字符h的状态转移公式，查表：

  - ```base[s]+c = base[116]+'h' = 4+104 = 108```
  - ```check[108] = 116```

  根据上面的公式，index[116]节点下找到与字符h相等的孩子节点index[108]，返回孩子节点的index给currentIndex，即**currentIndex =108**。

- 查表显示lengths[108]为空，说明当前位置没有结束字符。

- 继续下一个字符查找。

![avatar](./images/search_eleventh.jpg)

###### 12、查找第12个字符e，ascii码是101

- 当前**currentIndex=108**。

- 查表显示base[108]不为空，说明index[108]下有孩子节点。

  查看index[108]节点下如果有字符e的状态转移公式，查表：

  - ```base[s]+c = base[108]+'e' = 1+101 = 102```
  - ```check[102] = 108```

  根据上面的公式可知，index[108]节点下找到与e相等的孩子节点index[102]，返回孩子节点的index给currentIndex，即**currentIndex =102**。

- 查表显示lengths[102]=3,2，说明当前位置有2个结束字符，第一个模式串长度是3，起始坐标(9,11)，匹配模式串：she；第二个模式串长度是2，起始坐标(10,11)，匹配到模式串：he。

- 继续下一个字符查找。

![avatar](./images/search_twelfth.jpg)

###### 13、查找第13个字符h，ascii码是104

- 当前**currentIndex=102**。

- 查表显示base[102]为空，说明index[102]节点下没有孩子节点，并且index[102]不是root节点，查表可知fail[102]=103，将103返回给currentIndex，即**currentIndex=103**。

- 查表显示base[103]不为空，说明index[103]节点下有孩子节点。

  查看index[103]节点下如果有字符h的状态转移公式，查表：

  - ```base[s]+c = base[103]+'h' = 1+104 = 105```
  - ```check[105] = 0```

  根据上面的公式可知：check[105] != 103。说明index[103]节点下没有与h相等的孩子节点，查表可知fail[103]=0，返回0给currentIndex，即**currentIndex=0**。

- 查表显示base[0]不为空，说明index[0]节点下有孩子节点。

- 查看index[0]节点下如果有字符h的状态转移公式，查表：

  - ```base[s]+c = base[0]+'h' = 1+104 = 105```
  - ```check[105] = 0```

  根据上面公式可知，index[0]节点下找到了与字符h相等的孩子节点index[105]，返回孩子节点的index给currentIndex，即**currentIndex=105**。

- 查表显示lengths[105]为空，说明当前位置没有结束字符。

- 继续下一个字符查找。

![avatar](./images/search_thirteenth.jpg)

###### 14、查找第14个字符e，ascii码是101

- 当前**currentIndex=105**。

- 查表显示base[105]不为空，说明index[105]下有孩子节点。

  查看index[105]节点下如果有字符e的状态转移公式，查表：

  - ```base[s]+c = base[105]+'e' = 2+101 = 103```
  - ```check[103] = 105```

  根据上面的公式可知，index[105]节点下找到与字符e相等的孩子节点index[103]，返回孩子节点的index给currentIndex，即**currentIndex=103**。

- 查表显示lengths[103]=2，说明当前位置有1个结束字符，这个模式串长为2，起始坐标(12,13)，匹配的模式串是he。

- 继续下一个字符查找。

![avatar](./images/search_fourteenth.jpg)

###### 15、查找第15个字符r，ascii码是114

- 当前**currentIndex=103**。

- 查表显示base[103]不为空，说明index[103]节点下有孩子节点。

  查看index[103]节点下如果有字符r的状态转移公式，查表：

  - ```base[s]+c = base[103]+'r' = 1+114 = 115```
  - ```check[115] = 103```

  根据上面公式可知，index[103]节点下找到与字符r相等的孩子节点index[115]，返回孩子节点的index给currentIndex，即**currentIndex=115**。

- 查表显示lengths[103]为空，说明当前位置没有结束字符。

- 继续下一个字符查找。

![avatar](./images/search_fifteenth.jpg)

###### 16、查找第16个字符s，ascii码是115

- 当前**currentIndex=115**。

- 查表显示base[115]不为空，说明index[115]有孩子节点。

  查看index[115]节点下如果有字符s的状态转移公式，查表：

  - ```base[s]+c = base[115]+'s' = 3+115 = 118```
  - ```check[118] = 115```

  根据上面公式可知，index[115]节点下找到与字符s相等的孩子节点index[118]，返回孩子节点的index给currentIndex，即**currentIndex=118**。

- 查表显示lengths[118]=4，说明当前位置有1个结束字符，这个模式串长4，起始坐标(12, 15)，匹配到的模式串是hers。

- 继续下一个字符查找。

![avatar](./images/search_sixteenth.jpg)

###### 17、接下来的流程不再赘述。

###### 三、最后的结果

![avatar](./images/result.jpg)

- 基于上面的查询流程可以知道，对于文本串"ifindhehishehersall"，一共19个字符，只需要19次字符查找，就能查找出所有这些模式串[i，he，his，she，hers]在给定文本串中出现的位置信息，查询次数跟模式串数量的多少没有关系，只跟被查询的文本串长度有关。
- 如果被比较的字符是汉字等可变长字符，使用golang的rune类型，也能实现一次查询是否存在。

## 总结

###### 优点：

- double array trie能精确匹配多个模式串，匹配次数跟给定的文本串长度相关，跟具体的模式串个数多少关系不大（模式串比较多会导致沿着fail数组的比较次数增多），比使用hash实现的trie的查询速度的理论上更快。
- 加入fail数组后，只需要扫描一次，便能找出文本串中所有出现的模式串。
- 在golang版本实现中，直接是以rune类型记录在base和check数组，也是以rune类型比较，一个汉字3个字节，这样操作以后汉字查找的比较次数会更少，速度会更快。

###### 缺点：

- 动态增删模式串需要调整的节点比较多，实现起来比较复杂。参考我们项目的屏蔽字检测功能，这个需求比较少。

###### 其他：

- 如果需要对找到的模式串过滤白名单，可以再额外建立一份白名单双数组字典树，最后得到2份查找结果，把符合的白名单模式串下标的屏蔽字过滤掉。

- double array trie在构造的base和check数组的时候还是相对耗时的，如果项目对时间消耗敏感，可以将base、check、fail和lengths本地序列化，然后在项目里加载已经构建好的相关数组；如果项目对时间消耗不敏感，可以单独起一个go协程在后台慢慢构建双数组字典树。

## 代码实现

[双数组字典树Golang实现](https://github.com/theflavoroflife/double_array_trie)

## 后记

这篇文章是叠纸恋与制作人程序组的技术交流与分享，转载请注明。

## 参考资料

[1] [Compressed double-array tries for string dictionaries supporting fast lookup](https://kampersanda.github.io/pdf/KAIS2017.pdf)

[2] [An Efficient Implementation of Trie Structures](https://www.co-ding.com/assets/pdf/dat.pdf) 

[3] [An Implementation of Double-Array Trie](https://linux.thai.net/~thep/datrie/datrie.html)

[4] [Aho–Corasick algorithm](https://oi-wiki.org/string/ac-automaton/)

