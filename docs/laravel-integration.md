# Laravel Integration with lazySMTP

This guide shows how to integrate lazySMTP into your Laravel application for email testing during development.

## Configuration

### Update Laravel's Mail Configuration

Edit your `.env` file:

```env
MAIL_MAILER=smtp
MAIL_HOST=localhost
MAIL_PORT=2525
MAIL_USERNAME=null
MAIL_PASSWORD=null
MAIL_ENCRYPTION=null
MAIL_FROM_ADDRESS="noreply@yourapp.com"
MAIL_FROM_NAME="${APP_NAME}"
```

## Sending Test Emails

### Basic Email

```php
use Illuminate\Support\Facades\Mail;

Mail::raw('This is a test email', function ($message) {
    $message->to('test@example.com')
            ->subject('Test Email from Laravel');
});
```

### Mailable Classes

Create a mailable:

```bash
php artisan make:mail TestEmail
```

In `app/Mail/TestEmail.php`:

```php
namespace App\Mail;

use Illuminate\Bus\Queueable;
use Illuminate\Mail\Mailable;
use Illuminate\Queue\SerializesModels;

class TestEmail extends Mailable
{
    use Queueable, SerializesModels;

    public $data;

    public function __construct($data = [])
    {
        $this->data = $data;
    }

    public function build()
    {
        return $this->view('emails.test')
                    ->subject('Test Email')
                    ->with($this->data);
    }
}
```

Send the mailable:

```php
use App\Mail\TestEmail;
use Illuminate\Support\Facades\Mail;

Mail::to('recipient@example.com')->send(new TestEmail([
    'name' => 'John Doe',
]));
```

### Sending with Attachments

```php
Mail::raw('Email body', function ($message) {
    $message->to('recipient@example.com')
            ->subject('Email with Attachment')
            ->attach(storage_path('app/document.pdf'));
});
```

### Sending with Blade Template

Create `resources/views/emails/test.blade.php`:

```blade
<h1>Hello, {{ $name }}</h1>

<p>This is a test email sent from Laravel.</p>

<p>Best regards,<br>The Team</p>
```

Send with the template:

```php
Mail::send('emails.test', ['name' => 'John'], function ($message) {
    $message->to('john@example.com')
            ->subject('Welcome!');
});
```

## Testing Email Functionality

### Using Laravel Tinker

```bash
php artisan tinker
```

```php
Mail::raw('Test email from Tinker', function ($message) {
    $message->to('test@example.com')
            ->subject('Tinker Test');
});
```

### In PHPUnit Tests

```php
use Illuminate\Support\Facades\Mail;

public function test_email_sending()
{
    Mail::fake();

    Mail::raw('Test', function ($message) {
        $message->to('test@example.com')
                ->subject('Test');
    });

    Mail::assertSent(function ($mail) {
        return $mail->hasTo('test@example.com');
    });
}
```

## Notification Channels

LazySMTP works with Laravel's notification system:

```php
namespace App\Notifications;

use Illuminate\Bus\Queueable;
use Illuminate\Notifications\Notification;
use Illuminate\Notifications\Messages\MailMessage;

class TestNotification extends Notification
{
    use Queueable;

    public function via($notifiable)
    {
        return ['mail'];
    }

    public function toMail($notifiable)
    {
        return (new MailMessage)
                    ->greeting('Hello!')
                    ->line('This is a test notification.')
                    ->action('View Action', url('/'))
                    ->line('Thank you for using our application!');
    }
}
```

Send the notification:

```php
use App\Notifications\TestNotification;
use App\Models\User;

$user = User::first();
$user->notify(new TestNotification());
```

## Common Use Cases

### Password Reset Emails

Laravel's built-in password reset will automatically use lazySMTP when configured:

```bash
php artisan tinker
```

```php
use App\Models\User;

$user = User::first();
$user->sendPasswordResetNotification('reset-token-here');
```

### Verification Emails

```php
$user->sendEmailVerificationNotification();
```

### Queue Email Jobs

```bash
php artisan queue:work
```

Queued emails will be sent to lazySMTP when processed.

## Troubleshooting

### Email Not Received

1. Ensure lazySMTP is running
2. Check the port configuration matches (default: 2525)
3. Verify Laravel's `.env` settings
4. Clear configuration cache: `php artisan config:clear`

### Connection Refused

- Make sure lazySMTP is started and showing "Running" status
- Check firewall settings aren't blocking port 2525

### Email Format Issues

- Use lazySMTP's email viewer to inspect raw email content
- Check that your mail views are properly formatted

## Production Deployment

Remember to update your `.env` file with production mail settings (e.g., SMTP, SendGrid, Mailgun, SES) before deploying:

```env
MAIL_MAILER=smtp
MAIL_HOST=smtp.your-provider.com
MAIL_PORT=587
MAIL_USERNAME=your-username
MAIL_PASSWORD=your-password
MAIL_ENCRYPTION=tls
```

## Tips

1. **View in lazySMTP TUI**: Start lazySMTP with `./lazysmtp` and monitor emails in real-time
2. **Persistent Database**: Emails are saved to `lazysmtp.db` and can be reviewed later
3. **Quick Reset**: Delete the database file to clear all emails
4. **Multiple Projects**: Use different database paths for different projects: `./lazysmtp -db project1.db`

## Resources

- [Laravel Mail Documentation](https://laravel.com/docs/mail)
- [Laravel Notifications](https://laravel.com/docs/notifications)
- [Laravel Queues](https://laravel.com/docs/queues)
