FROM centos:7

# 安装必要的软件
RUN yum install -y wget tar java-1.8.0-openjdk-devel

# 下载并解压 Hadoop
RUN wget http://mirror.bit.edu.cn/apache/hadoop/common/hadoop-3.3.1/hadoop-3.3.1.tar.gz && \
    tar -zxvf hadoop-3.3.1.tar.gz && \
    mv hadoop-3.3.1 /usr/local/hadoop

# 设置环境变量
ENV JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk
ENV HADOOP_HOME=/usr/local/hadoop
ENV PATH=$PATH:$JAVA_HOME/bin:$HADOOP_HOME/bin:$HADOOP_HOME/sbin

# 格式化 HDFS
RUN hdfs namenode -format

# 启动 Hadoop
CMD ["/usr/local/hadoop/sbin/start-all.sh"]