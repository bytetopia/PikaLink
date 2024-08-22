from django.urls import path

from . import views

urlpatterns = [
    path('<path:full_path>/', views.jump, name='jump'),
    path('', views.index_page, name='index_page')
]