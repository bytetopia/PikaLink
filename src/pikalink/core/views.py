from django.conf import settings
from django.shortcuts import render
from django.http import HttpResponse, HttpResponseRedirect
from . import models


def jump(req, full_path):
    # get first part of path
    path_parts = full_path.split('/')
    short_path = path_parts[0]
    remaining_path = '/'.join(path_parts[1:]) if len(path_parts) > 1 else ''

    # get request parameters
    query_string = req.META.get('QUERY_STRING', '')
    user_agent = req.META.get('HTTP_USER_AGENT', '')
    user_referer = req.META.get('HTTP_REFERER', '')

    # special logic to fetch real ip
    x_forwarded_for = req.META.get('HTTP_X_FORWARDED_FOR')
    if x_forwarded_for:
        user_ip = x_forwarded_for.split(',')[0]
    else:
        user_ip = req.META.get('REMOTE_ADDR', '')

    # prepare return context
    return_context = {
        'title': '',
        'msg': '',
        'site_title': settings.SITE_NAME,
    }

    # check is short link exist
    links = models.ShortUrl.objects.filter(short_path=short_path, is_deleted=0)
    if links.count() == 1:
        link = links[0]

        # check auto invalid
        if link.is_auto_invalid == 1:
            # check if invalid
            if link.invalid_after_time < datetime.datetime.now():
                link.is_deleted = 1
                link.save()
                return_context['title'] = 'Invalid Link'
                return_context['msg'] = 'The short link is expired.'
                return render(req, 'core/index.html', return_context)

        # write log
        record = models.AccessRecord()
        record.url = link
        record.source_ip = user_ip
        record.referer = user_referer
        record.user_agent = user_agent
        record.save()

        # redirect
        redir_path = link.full_path
        if remaining_path != '':
            redir_path += '/' + remaining_path
        if query_string != '':
            redir_path += '?' + query_string
        return HttpResponseRedirect(redir_path)
    # link does not exist
    else:
        return_context['title'] = 'Invalid Link'
        return_context['msg'] = 'Link does not exist or expired.'
        return render(req, 'index.html', return_context)


def index_page(req):
    context = {
        'title': settings.SITE_NAME,
        'msg': settings.SITE_DESCRIPTION,
        'site_title': settings.SITE_NAME,
    }
    return render(req, 'index.html', context=context)