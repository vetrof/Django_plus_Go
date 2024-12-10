from rest_framework import serializers
from rest_framework.views import APIView
from rest_framework.response import Response
import requests
from django.http import JsonResponse

from catalog.models import Catalog


# Сериализатор для продукта
class ProductSerializer(serializers.ModelSerializer):
    class Meta:
        model = Catalog
        fields = ['id', 'name', 'description', 'price']  # Указываем только нужные поля


# Пример использования сериализатора
class Products(APIView):
    def get(self, request):
        products = Catalog.objects.all()

        # Сериализация данных
        serializer = ProductSerializer(products, many=True)
        return Response(serializer.data)



def g_products(request):
    # Отправка GET-запроса в Go-сервис
    response = requests.get('http://go_service:8080/product')  # Пример URL для Go-сервиса
    # Возвращаем ответ в формате JSON
    return JsonResponse(response.json(), safe=False)
