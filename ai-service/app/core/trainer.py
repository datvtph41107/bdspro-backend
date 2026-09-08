"""
Model Trainer - Fine-tune sentence transformers on custom Q&A data
"""

import json
from pathlib import Path
from typing import List, Dict, Tuple
from sentence_transformers import SentenceTransformer, InputExample, losses
from torch.utils.data import DataLoader
import logging

logger = logging.getLogger(__name__)

class ModelTrainer:
    """Handles model training/fine-tuning"""
    
    def __init__(
        self,
        base_model: str = "all-MiniLM-L6-v2",
        output_path: str = "./models/fine_tuned"
    ):
        """
        Initialize trainer
        
        Args:
            base_model: Base model to fine-tune
            output_path: Path to save fine-tuned model
        """
        self.base_model = base_model
        self.output_path = Path(output_path)
        self.output_path.mkdir(parents=True, exist_ok=True)
        
        logger.info(f"Initializing trainer with base model: {base_model}")
        self.model = None
    
    def load_training_data(self, data_path: str) -> List[InputExample]:
        """
        Load Q&A pairs from JSON file(s)
        
        Args:
            data_path: Path to a JSON file or directory containing JSON files
            
        Returns:
            List of InputExample for training
        """
        logger.info(f"Loading training data from: {data_path}")
        
        path = Path(data_path)
        qa_pairs = []
        
        # Check if it's a directory or a file
        if path.is_dir():
            # Load all JSON files from directory
            json_files = list(path.glob("*.json"))
            if not json_files:
                raise ValueError(f"No JSON files found in directory: {data_path}")
            
            logger.info(f"Found {len(json_files)} JSON files in directory")
            
            for json_file in json_files:
                logger.info(f"  Loading: {json_file.name}")
                with open(json_file, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                    if isinstance(data, list):
                        qa_pairs.extend(data)
                    else:
                        logger.warning(f"Skipping {json_file.name}: not a JSON array")
        else:
            # Load single file
            logger.info(f"Loading single file: {path.name}")
            with open(path, 'r', encoding='utf-8') as f:
                data = json.load(f)
                if isinstance(data, list):
                    qa_pairs = data
                else:
                    raise ValueError(f"Invalid format: {data_path} must contain a JSON array")
        
        # Convert to InputExample format
        # For sentence transformers, we create pairs of (question, answer)
        examples = []
        
        for idx, pair in enumerate(qa_pairs):
            question = pair.get('question', '')
            answer = pair.get('answer', '')
            
            if not question or not answer:
                logger.warning(f"Skipping invalid pair at index {idx}: missing question or answer")
                continue
            
            # Create example with question and answer as a pair
            # Label 1.0 means they are semantically similar
            examples.append(
                InputExample(texts=[question, answer], label=1.0)
            )
        
        logger.info(f"Loaded {len(examples)} training examples from {len(qa_pairs)} Q&A pairs")
        return examples
    
    def train(
        self,
        train_data: List[InputExample],
        epochs: int = 3,
        batch_size: int = 16,
        warmup_steps: int = 100
    ) -> Dict:
        """
        Fine-tune the model
        
        Args:
            train_data: List of training examples
            epochs: Number of training epochs
            batch_size: Batch size
            warmup_steps: Warmup steps for learning rate
            
        Returns:
            Training statistics
        """
        logger.info(f"Starting training with {len(train_data)} examples")
        logger.info(f"Epochs: {epochs}, Batch size: {batch_size}")
        
        # Load base model
        logger.info(f"Loading base model: {self.base_model}")
        self.model = SentenceTransformer(self.base_model)
        
        # Create DataLoader
        train_dataloader = DataLoader(
            train_data,
            shuffle=True,
            batch_size=batch_size
        )
        
        # Define loss function
        # CosineSimilarityLoss is good for Q&A matching
        train_loss = losses.CosineSimilarityLoss(self.model)
        
        # Calculate training steps
        total_steps = len(train_dataloader) * epochs
        
        logger.info(f"Training for {total_steps} steps...")
        
        # Train the model
        self.model.fit(
            train_objectives=[(train_dataloader, train_loss)],
            epochs=epochs,
            warmup_steps=warmup_steps,
            output_path=str(self.output_path),
            show_progress_bar=True,
            save_best_model=True
        )
        
        logger.info(f"Training complete! Model saved to: {self.output_path}")
        
        return {
            "status": "success",
            "model_path": str(self.output_path),
            "epochs": epochs,
            "training_examples": len(train_data),
            "total_steps": total_steps
        }
    
    def evaluate(self, test_data: List[InputExample]) -> Dict:
        """
        Evaluate the fine-tuned model
        
        Args:
            test_data: List of test examples
            
        Returns:
            Evaluation metrics
        """
        if self.model is None:
            self.model = SentenceTransformer(str(self.output_path))
        
        logger.info(f"Evaluating on {len(test_data)} examples")
        
        # Simple evaluation: check if question embeddings are close to answer embeddings
        correct = 0
        total = len(test_data)
        
        for example in test_data:
            q_embedding = self.model.encode(example.texts[0])
            a_embedding = self.model.encode(example.texts[1])
            
            # Calculate cosine similarity
            from numpy import dot
            from numpy.linalg import norm
            
            similarity = dot(q_embedding, a_embedding) / (norm(q_embedding) * norm(a_embedding))
            
            # If similarity > threshold, consider correct
            if similarity > 0.5:
                correct += 1
        
        accuracy = correct / total if total > 0 else 0
        
        logger.info(f"Evaluation complete! Accuracy: {accuracy:.2%}")
        
        return {
            "accuracy": accuracy,
            "correct": correct,
            "total": total
        }

def train_model(
    data_path: str = "./data/train",
    output_path: str = "./models/fine_tuned",
    epochs: int = 3,
    batch_size: int = 16
) -> Dict:
    """
    Convenience function to train model
    
    Args:
        data_path: Path to training data (file or directory with JSON files)
        output_path: Path to save model
        epochs: Number of epochs
        batch_size: Batch size
        
    Returns:
        Training results
    """
    trainer = ModelTrainer(output_path=output_path)
    
    # Load data (will read all JSON files if data_path is a directory)
    train_data = trainer.load_training_data(data_path)
    
    # Train
    results = trainer.train(
        train_data=train_data,
        epochs=epochs,
        batch_size=batch_size
    )
    
    return results

