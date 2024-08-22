from django.db import models

class ShortUrl(models.Model):
    YES_NO_CHOICES = (
        (0, 'No'),
        (1, 'Yes')
    )
    # Link information
    short_path = models.CharField(max_length=30, null=False, unique=True, help_text='URL slug')
    full_path = models.TextField(default='', null=False)
    link_desc = models.CharField(max_length=100, blank=True, default='', help_text='link description, not visible to user')
    # Auto invalid
    is_auto_invalid = models.IntegerField(choices=YES_NO_CHOICES, default=0, help_text='link will expire after given time')
    invalid_after_time = models.DateTimeField(blank=True, help_text='invalid after this time', null=True)
    # Other fields
    create_time = models.DateTimeField(auto_now_add=True, help_text='create time', blank=True, null=True)
    update_time = models.DateTimeField(auto_now=True, help_text='update time', blank=True, null=True)
    is_deleted = models.IntegerField(choices=YES_NO_CHOICES, default=0, help_text='is marked as deleted')

    def __str__(self):
        return self.short_path


class AccessRecord(models.Model):
    YES_NO_CHOICES = (
        (0, 'No'),
        (1, 'Yes')
    )
    url = models.ForeignKey(ShortUrl, on_delete=models.CASCADE, help_text='Short url')
    time = models.DateTimeField(auto_now_add=True, help_text='Visit time')
    source_ip = models.CharField(max_length=50)
    source_ip_country = models.CharField(max_length=100)
    source_ip_city = models.CharField(max_length=100)
    referer = models.TextField()
    referer_domain = models.CharField(max_length=100)
    user_agent = models.TextField()
    user_agent_device = models.CharField(max_length=200)
    user_agent_os = models.CharField(max_length=200)
    user_agent_browser = models.CharField(max_length=200)

    def __str__(self):
        return '%s - %s - %s' % (self.url.short_path, self.time.strftime('%Y-%m-%d %H:%M:%S'), self.source_ip)