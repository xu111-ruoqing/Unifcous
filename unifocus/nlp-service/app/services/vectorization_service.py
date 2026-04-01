from sentence_transformers import SentenceTransformer
from loguru import logger
import os

class VectorizationService:
    def __init__(self):
        self.model_name = os.getenv("NLP_MODEL_NAME", "all-MiniLM-L6-v2")
        self.model = None

    def load_model(self):
        """Lazy load the model"""
        if self.model is None:
            logger.info(f"Loading SentenceTransformer model: {self.model_name}")
            try:
                # 尝试从本地加载，如果不存在则自动下载
                # 在Docker环境中，建议挂载模型目录到 /root/.cache/torch/sentence_transformers
                self.model = SentenceTransformer(self.model_name)
                logger.info("Model loaded successfully")
            except Exception as e:
                logger.error(f"Failed to load model: {e}")
                raise e

    def vectorize(self, text: str) -> list[float]:
        """Convert text to vector"""
        self.load_model()
        try:
            # model.encode returns numpy array, convert to list
            embeddings = self.model.encode(text)
            return embeddings.tolist()
        except Exception as e:
            logger.error(f"Vectorization failed: {e}")
            raise e

# Global instance
vectorization_service = VectorizationService()
