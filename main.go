// Package main はアプリケーションのエントリーポイントです。
// このパッケージは、サーバーの起動とルーティングの初期化を行います。
package main

import (
	"golang-clean-architecture-example/infrastructure"
)

// main はアプリケーションのメイン関数です。
// ルーターを初期化し、Echo サーバーを起動します。
func main() {
	infrastructure.InitRouter()
}
