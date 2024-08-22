from django.contrib import admin
from unfold.admin import ModelAdmin

from .models import ShortUrl, AccessRecord


@admin.register(ShortUrl)
class ShortUrlAdminClass(ModelAdmin):
    pass

@admin.register(AccessRecord)
class AccessRecordAdminClass(ModelAdmin):
    pass