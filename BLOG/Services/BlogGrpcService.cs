using Blog; // This namespace comes from the Proto generation
using Grpc.Core;
using BLOG.Model;
using BLOG.Services;
using Google.Protobuf.WellKnownTypes;
using System.Linq;
using System.Runtime.CompilerServices;

namespace BLOG.GrpcServices
{
    public class BlogGrpcService : Blog.BlogService.BlogServiceBase
    {
        private readonly Services.BlogService _blogService;
        private readonly CommentService _commentService;

        public BlogGrpcService(Services.BlogService blogService, CommentService commentService)
        {
            _blogService = blogService;
            _commentService = commentService;
        }

        // Helper to get UserID from Metadata (sent by Go Gateway)
        private string GetUserName(ServerCallContext context)
        {
            var userEntry = context.RequestHeaders.FirstOrDefault(h => h.Key == "user-username");

            if (userEntry == null || string.IsNullOrEmpty(userEntry.Value))
            {
                throw new RpcException(new Status(StatusCode.Unauthenticated, "No userName found"));
            }
            return userEntry.Value;
        }

        public override async Task<GetAllBlogsResponse> GetAllPosts(Empty request, ServerCallContext context)
        {
            var blogs = await _blogService.GetBlogsAsync();
            var response = new GetAllBlogsResponse();

            foreach (var blog in blogs)
            {
                var blogt = new Blog.Blog
                {
                    Id = blog.Id,
                    Title = blog.Title,
                    Description = blog.Description,
                    UserName = blog.UserName,
                    LikeCount = blog.LikeCount,
                    CreatedAt = Timestamp.FromDateTime(blog.CreatedAt.ToUniversalTime())
                };
                response.Blogs.Add(blogt);
            }

            return response;
        }

        public override async Task<BlogResponse> CreatePost(CreateBlogRequest request, ServerCallContext context)
        {
            var userName = GetUserName(context);
            var newBlog = new Model.Blog
            {
                Title = request.Blog.Title,
                UserName = userName,
                Description = request.Blog.Description,
                ImagePaths = request.Blog.ImagePaths.ToList(),
                CreatedAt = DateTime.UtcNow,
                LikeCount = 0
            };

            await _blogService.CreateBlogAsync(newBlog);
            return MapToProtoBlog(newBlog);
        }

        public override async Task<BlogResponse> ToggleLike(ToggleLikeRequest request, ServerCallContext context)
        {
            var userName = GetUserName(context);
            var updatedBlog = await _blogService.ToggleBlogLikeAsync(request.PostId, userName);

            if (updatedBlog == null)
                throw new RpcException(new Status(StatusCode.NotFound, "Blog not found"));

            return MapToProtoBlog(updatedBlog);
        }

        // ======================
        // COMMENT METHODS
        // ======================
        public override async Task<GetCommentsResponse> GetComments(GetCommentsRequest request, ServerCallContext context)
        {
            var comments = await _commentService.GetCommentsByPostIdAsync(request.PostId);
            var response = new GetCommentsResponse();

            foreach (var comment in comments)
            {
                var newComment = new Blog.Comment
                {
                  Id = comment.Id,
                  PostId = request.PostId,
                  Text = comment.Text,
                  AuthorName = comment.AuthorName,
                  CreatedAt = Timestamp.FromDateTime(comment.CreatedAt.ToUniversalTime()),
                  UpdatedAt = Timestamp.FromDateTime(comment.UpdatedAt.ToUniversalTime())
                };
                response.Comments.Add(newComment);
            }

            return response;
        }

        public override async Task<CommentResponse> CreateComment(CreateCommentRequest request, ServerCallContext context)
        {
            var userName = GetUserName(context);
            var newComment = new Model.Comment
            {
                Text = request.Comment.Text,
                AuthorName = userName,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _commentService.CreateCommentAsync(newComment);
            return MapToProtoComment(newComment);
        }

        public override async Task<CommentResponse> UpdateComment(UpdateCommentRequest request, ServerCallContext context)
        {
            var userName = GetUserName(context);
            var newComment = new Model.Comment
            {
                Id = request.Comment.Id,
                PostId = request.Comment.PostId,
                Text = request.Comment.Text,
                AuthorName = userName,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            var success = await _commentService.UpdateCommentAsync(newComment);

            if (!success)
                throw new RpcException(new Status(StatusCode.NotFound, "Comment not found"));
            return MapToProtoComment(newComment);
        }

        public override async Task<DeleteCommentResponse> DeleteComment(DeleteCommentRequest request, ServerCallContext context)
        {
            var userName = GetUserName(context);
            await _commentService.DeleteCommentAsync(request.CommentId);
            return new DeleteCommentResponse { Success = true };
        }

        // ======================
        // Mappers: Domain -> Proto
        // ======================
        private static BlogResponse MapToProtoBlog(Model.Blog blog)
        {
            var resp = new BlogResponse
            {
                Blog = new Blog.Blog
                {
                    Id = blog.Id,
                    Title = blog.Title,
                    Description = blog.Description,
                    UserName = blog.UserName,
                    LikeCount = blog.LikeCount,
                    CreatedAt = Timestamp.FromDateTime(blog.CreatedAt.ToUniversalTime())
                }
            };

            if (blog.ImagePaths != null)
                resp.Blog.ImagePaths.AddRange(blog.ImagePaths);

            return resp;
        }

        private static CommentResponse MapToProtoComment(Model.Comment comment)
        {
            return new CommentResponse
            {
                Comment = new Blog.Comment
                {
                    Id = comment.Id,
                    PostId = comment.PostId,
                    AuthorName = comment.AuthorName,
                    Text = comment.Text,
                    CreatedAt = Timestamp.FromDateTime(comment.CreatedAt.ToUniversalTime()),
                    UpdatedAt = Timestamp.FromDateTime(comment.UpdatedAt.ToUniversalTime())
                }
            };
        }
    }
}