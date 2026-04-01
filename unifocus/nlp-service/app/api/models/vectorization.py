from pydantic import BaseModel
from typing import List

class VectorizeRequest(BaseModel):
    text: str

class VectorizeResponse(BaseModel):
    vector: List[float]
    dim: int
