"""
文本提取API路由
"""
from fastapi import APIRouter, UploadFile, File, HTTPException
from pydantic import BaseModel
from loguru import logger

from app.services.text_extractor import TextExtractor

router = APIRouter()
extractor = TextExtractor()


class HTMLRequest(BaseModel):
    """HTML提取请求"""
    html: str


class TextResponse(BaseModel):
    """文本提取响应"""
    text: str
    length: int


@router.post("/html", response_model=TextResponse, tags=["Text Extraction"])
async def extract_html(request: HTMLRequest):
    """
    从HTML中提取纯文本
    
    - **html**: HTML内容字符串
    """
    try:
        text = extractor.extract_from_html(request.html)
        return TextResponse(text=text, length=len(text))
    except Exception as e:
        logger.error(f"HTML extraction error: {e}")
        raise HTTPException(status_code=400, detail=str(e))


@router.post("/pdf", response_model=TextResponse, tags=["Text Extraction"])
async def extract_pdf(file: UploadFile = File(...)):
    """
    从PDF文件中提取文本
    
    - **file**: PDF文件（支持文件上传）
    """
    try:
        # 验证文件类型
        if not file.filename.endswith('.pdf'):
            raise HTTPException(status_code=400, detail="File must be a PDF")
        
        # 读取文件内容
        pdf_bytes = await file.read()
        
        # 提取文本
        text = extractor.extract_from_pdf(pdf_bytes)
        
        return TextResponse(text=text, length=len(text))
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"PDF extraction error: {e}")
        raise HTTPException(status_code=400, detail=f"PDF extraction failed: {str(e)}")


@router.post("/docx", response_model=TextResponse, tags=["Text Extraction"])
async def extract_docx(file: UploadFile = File(...)):
    """
    从Word文档中提取文本
    
    - **file**: Word文档文件（.docx格式）
    """
    try:
        # 验证文件类型
        if not (file.filename.endswith('.docx') or file.filename.endswith('.doc')):
            raise HTTPException(status_code=400, detail="File must be a DOCX document")
        
        # 读取文件内容
        docx_bytes = await file.read()
        
        # 提取文本
        text = extractor.extract_from_docx(docx_bytes)
        
        return TextResponse(text=text, length=len(text))
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"DOCX extraction error: {e}")
        raise HTTPException(status_code=400, detail=f"DOCX extraction failed: {str(e)}")

