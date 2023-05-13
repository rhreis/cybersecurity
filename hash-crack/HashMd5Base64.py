import hashlib
import base64

# Nome do arquivo de entrada com as senhas
arquivo_senhas = "/usr/share/john/password.lst"

# Nome do arquivo de saída para salvar os hashes em base64
arquivo_hash_base64 = "hash_base64.txt"

# Abrir o arquivo de entrada com as senhas
with open(arquivo_senhas, "r") as arquivo:
    # Ler as senhas uma por uma
    senhas = arquivo.readlines()

    # Abrir o arquivo de saída para salvar os hashes em base64
    with open(arquivo_hash_base64, "w") as arquivo_saida:
        # Iterar sobre cada senha
        for senha in senhas:
            # Remover o caractere de quebra de linha
            senha = senha.strip()

            # Aplicar o algoritmo MD5 na senha
            md5 = hashlib.md5()
            md5.update(senha.encode("utf-8"))
            hash_md5 = md5.hexdigest()

            # Codificar o resultado do MD5 em base64
            hash_base64 = base64.b64encode(hash_md5.encode("utf-8")).decode("utf-8")

            # Codificar o resultado do base64 em sah1
            sha1 = hashlib.sha1()
            sha1.update(hash_base64.encode("utf-8"))
            hash_sha1 = sha1.hexdigest()

            # Salvar o hash em base64 no arquivo de saída
            arquivo_saida.write(senha + " - " + hash_sha1 + "\n")

# Exibir mensagem de conclusão
print("Hashes em base64 foram salvos em", arquivo_hash_base64)
