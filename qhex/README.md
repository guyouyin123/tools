[toc]

# 编码：base64、URL、ASCII、UTF-8、UTF-16、hex

# 更多加密参考：https://github.com/duke-git/lancet/blob/main/docs/api/packages/cryptor.md#AesEcbEncrypt
```go
1.Base64: 用于将二进制数据编码为 ASCII 字符串，便于文本传输和存储。

2.URL 编码: 将 URL 中的特殊字符转化为百分号（%）后跟两个十六进制数字的形式，以便在 URL 中安全传输。

3.ASCII: 7 位编码标准，表示英文字符和一些控制字符，常用于基本文本处理。

4.UTF-8: 可变长度的编码方式，用于表示 Unicode 字符集中的字符，兼容 ASCII，并支持全球各种文字。

5.UTF-16: 16 位编码方式，也用于表示 Unicode 字符集，支持更多字符，但占用空间较大。

6.Hex 编码: 将每个字节的二进制数据转换为两个十六进制字符，常用于表示二进制数据的文本表示形式。
```

# 对称加密：AES、DES、3DES、Blowfish、RC4
对称加密算法使用相同的密钥进行加密和解密。常见的对称加密算法包括:

```go
1.AES (Advanced Encryption Standard):
密钥长度: 128, 192, 或 256 位
块大小: 128 位
用途: 广泛应用于数据加密，例如文件加密、通信加密等。

2.DES (Data Encryption Standard):
密钥长度: 56 位
块大小: 64 位
用途: 曾经广泛使用，但由于密钥长度较短，现已不推荐使用，主要被 AES 取代。

3.3DES (Triple DES):
密钥长度: 168 位（使用三个 56 位的密钥）
块大小: 64 位
用途: 提供比 DES 更高的安全性，但速度较慢，已逐步被 AES 取代。

4.Blowfish:
密钥长度: 32 至 448 位
块大小: 64 位
用途: 主要用于加密和解密，适用于多种应用。

5.RC4:
密钥长度: 40 至 2048 位（通常使用 128 位）
用途: 流加密算法，曾用于 SSL/TLS 和 WEP，但已因安全问题逐步弃用。
```

# 非对称加密：RSA、ECC、DSA、ElGamal

非对称加密算法使用一对密钥：公钥和私钥。常见的非对称加密算法包括：

```go
1.RSA (Rivest-Shamir-Adleman):
密钥长度: 通常为 1024, 2048, 或 4096 位
用途: 用于数据加密、数字签名和密钥交换等。

2.ECC (Elliptic Curve Cryptography):
密钥长度: 相对较短，提供相当的安全性（例如，256 位密钥与 3072 位 RSA 密钥的安全性相当）
用途: 用于加密和数字签名，越来越受欢迎，应用于 SSL/TLS 和各种安全协议。

3.DSA (Digital Signature Algorithm):
密钥长度: 通常为 1024 位或更长
用途: 主要用于数字签名，确保消息的完整性和认证。

4.ElGamal:
密钥长度: 通常为 2048 位或更长
用途: 用于加密和数字签名，但效率较低。
```

# 散列函数：MD5、SHA1、SHA256

虽然散列算法不是加密算法，它们在数据完整性检查中非常重要：

```go
1.MD5 (Message Digest Algorithm 5):
输出长度: 128 位
用途: 广泛用于文件校验和，但因安全漏洞不再推荐用于安全敏感应用。

2.SHA-1 (Secure Hash Algorithm 1):
输出长度: 160 位
用途: 主要用于数据完整性校验，已被发现存在安全漏洞，建议使用更强的哈希函数，如 SHA-256。

3.SHA-256 (Secure Hash Algorithm 256-bit):
输出长度: 256 位
用途: 广泛用于数据完整性校验和数字签名，提供较强的安全性
```

