from django.contrib import admin

from catalog.models import Catalog

@admin.register(Catalog)
class CataloAdmin(admin.ModelAdmin):
    list_display = ("name", "price")
    

