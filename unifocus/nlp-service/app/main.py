"""
UniFocus NLP Service - 主入口文件
提供文本提取、实体识别、OCR、向量化等NLP功能
"""
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from loguru import logger
import sys
from datetime import datetime

# 配置日志
logger.remove()
logger.add(
    sys.stdout,
    format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>",
    level="INFO"
)

# 创建 FastAPI 应用
app = FastAPI(
    title="UniFocus NLP Service",
    description="提供文本提取、实体识别、OCR、向量化等NLP功能",
    version="1.0.0"
)

# 配置 CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 启动事件
@app.on_event("startup")
async def startup_event():
    logger.info("Starting UniFocus NLP Service...")
    # TODO: 加载NLP模型
    # - Spacy中文模型
    # - Sentence Transformer模型
    # - PaddleOCR模型
    logger.info("NLP Service started successfully")

# 关闭事件
@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Shutting down NLP Service...")

# 健康检查
@app.get("/health")
async def health_check():
    return {
        "status": "ok",
        "service": "nlp-service",
        "version": "1.0.0",
        "timestamp": datetime.now().isoformat()
    }

# 根路径
@app.get("/")
async def root():
    return {
        "message": "UniFocus NLP Service API",
        "docs": "/docs",
        "health": "/health"
    }

# 引入路由
from app.api.routes import text_extractor
from app.api.routes import vectorization

app.include_router(text_extractor.router, prefix="/api/v1/extract", tags=["Text Extraction"])
app.include_router(vectorization.router, prefix="/api/v1/vectorize", tags=["Vectorization"])

# TODO: 其他路由待实现
# from app.api.routes import entity_recognizer, ocr_service
# app.include_router(entity_recognizer.router, prefix="/api/v1/entity", tags=["Entity Recognition"])
# app.include_router(ocr_service.router, prefix="/api/v1/ocr", tags=["OCR"])

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8000,
        reload=True,
        log_level="info"
    )
