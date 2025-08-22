"""
WohnFair Machine Learning Service

A comprehensive ML service for housing allocation fairness prediction,
demand forecasting, and risk assessment using advanced statistical models
and machine learning algorithms.
"""

__version__ = "0.1.0"
__author__ = "WohnFair Team"
__email__ = "team@wohnfair.de"

from . import models, preprocessing, training, evaluation, utils
from .config import get_settings

__all__ = [
    "__version__",
    "__author__",
    "__email__",
    "models",
    "preprocessing", 
    "training",
    "evaluation",
    "utils",
    "get_settings",
]
