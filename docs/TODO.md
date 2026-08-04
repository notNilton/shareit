# 📋 ShareIt Monorepo Roadmap & TODOs

Roadmap de desenvolvimento do **ShareIt** (Plataforma de compartilhamento de fotos com Go, React e React Native).

---

## 🚀 1. Backend (Go + MinIO + Redis + PostgreSQL)
- [ ] **Processamento Assíncrono de Imagens**
  - Implementar fila via Redis/Asynq para redimensionamento de thumbnails e geração de WebP em segundo plano.
- [ ] **Suporte a Transcodificação de Vídeo Curto**
  - Adicionar suporte a upload e conversão de pequenos vídeos via FFmpeg.

---

## 📱 2. Mobile & Web (React Native & React)
- [ ] **Cache Inteligente de Imagens Offline**
  - Implementar armazenamento em cache local de feed de fotos no mobile.
