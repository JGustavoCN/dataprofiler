package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

// A diretiva abaixo diz ao Go: "Inclua a pasta 'dist' inteira dentro do binário".
// Como este arquivo está DENTRO da pasta frontend, o caminho relativo é apenas "dist".
//
//go:embed dist/*
var distFS embed.FS

// GetFileSystem retorna o sistema de arquivos pronto para o servidor web.
// Ele já faz o "Sub" para entrar na pasta dist, assim o servidor acha o index.html na raiz.
func GetFileSystem() (http.FileSystem, error) {
	// Entra na subpasta "dist" do sistema de arquivos embutido
	fsys, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, err
	}
	return http.FS(fsys), nil
}
