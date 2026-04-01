from fastapi import APIRouter, HTTPException, Depends
from app.api.models.vectorization import VectorizeRequest, VectorizeResponse
from app.services.vectorization_service import vectorization_service

router = APIRouter()

@router.post("/text", response_model=VectorizeResponse)
async def vectorize_text(request: VectorizeRequest):
    """
    接收文本并返回向量
    """
    if not request.text.strip():
        raise HTTPException(status_code=400, detail="Text cannot be empty")
        
    try:
        vector = vectorization_service.vectorize(request.text)
        return VectorizeResponse(vector=vector, dim=len(vector))
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
